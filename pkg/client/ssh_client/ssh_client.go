package ssh_client

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/uuid"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"
)

type Client struct {
	SshClient  *ssh.Client
	SftpClient *sftp.Client
	Config     Config
	Stdout     io.Writer
	Stderr     io.Writer
}

type Config struct {
	Hostname     string   `yaml:"hostname"`
	Username     string   `yaml:"username"`
	Password     string   `yaml:"password"`
	PrivateKey   string   `yaml:"-"`
	AuthType     string   `yaml:"-"`
	Role         string   `yaml:"-"`
	Ip           string   `yaml:"ip"`
	Port         int      `yaml:"port"`
	KeyExchanges []string `yaml:"-"`
	Timeout      int      `yaml:"-"`
}

const (
	AuthTypePassword = "password"
	AuthTypeKey      = "key"
)

const (
	CONNECTION_TIMEOUT = 15 * time.Second
	TEXT_MAX_FILE_SIZE = 3 * 1024 * 1024
)

const (
	TypeDirectory = 1
	TypeFile      = iota
)

type OptionFunc func(c *Client)

func WithStdout(stdout io.Writer) OptionFunc {
	return func(c *Client) {
		c.Stdout = stdout
	}
}

func WithStderr(stderr io.Writer) OptionFunc {
	return func(c *Client) {
		c.Stderr = stderr
	}
}

func NewClient(config Config, optFunc ...OptionFunc) (*Client, error) {
	c := &Client{
		Config: config,
	}
	for _, of := range optFunc {
		of(c)
	}

	if err := c.connect(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) connect() error {

	if c.SshClient != nil {
		_, _, err := c.SshClient.SendRequest("keepalive", false, nil)
		if err == nil {
			return nil
		}
	}

	var auth ssh.AuthMethod

	if c.Config.AuthType == "" {
		auth = ssh.Password(c.Config.Password)
	}

	switch c.Config.AuthType {
	case AuthTypePassword:
		auth = ssh.Password(c.Config.Password)
	case AuthTypeKey:
		b, err := base64.StdEncoding.DecodeString(c.Config.PrivateKey)
		if err != nil {
			return err
		}
		signer, err := ssh.ParsePrivateKey(b)
		if err != nil {
			return fmt.Errorf("ssh parse private key: %w", err)
		}
		auth = ssh.PublicKeys(signer)
	}

	cfg := &ssh.ClientConfig{
		User: c.Config.Username,
		Auth: []ssh.AuthMethod{
			auth,
		},
		HostKeyCallback: func(string, net.Addr, ssh.PublicKey) error { return nil },
		Timeout:         CONNECTION_TIMEOUT,
		Config: ssh.Config{
			KeyExchanges: c.Config.KeyExchanges,
		},
	}

	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", c.Config.Ip, c.Config.Port), cfg)
	if err != nil {
		return fmt.Errorf("ssh dial: %w", err)
	}
	c.SshClient = sshClient

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return fmt.Errorf("sftp new client: %w", err)
	}
	c.SftpClient = sftpClient

	return nil
}

func (c *Client) Command(command string) (output string, err error) {

	if err := c.connect(); err != nil {
		return "", err
	}

	session, err := c.SshClient.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	if c.Config.Timeout <= 10 {

		var combinedOutput []byte
		combinedOutput, err = session.CombinedOutput(command)
		output = string(combinedOutput)
	} else {
		timeout := time.Duration(c.Config.Timeout) * time.Second

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		go func() {
			err = session.Run(command)
			if err != nil {
				log.Println("host:" + c.Config.Ip + "exec error:" + err.Error() + ",command:" + command)
			}
			cancel()
		}()

		select {
		case <-ctx.Done():
		case <-time.After(timeout):
			return "", errors.New("exec command timeout")
		}

		stdoutBuffer := new(bytes.Buffer)
		stderrBuffer := new(bytes.Buffer)

		stdoutWrites := make([]io.Writer, 0)
		stderrWrites := make([]io.Writer, 0)

		stdoutWrites = append(stdoutWrites, io.Writer(stdoutBuffer))
		stderrWrites = append(stderrWrites, io.Writer(stderrBuffer))

		if c.Stdout != nil {
			stdoutWrites = append(stdoutWrites, c.Stdout)
		}

		if c.Stderr != nil {
			stderrWrites = append(stderrWrites, c.Stderr)
		}

		session.Stdout = io.MultiWriter(stdoutWrites...)
		session.Stderr = io.MultiWriter(stderrWrites...)

		output = stdoutBuffer.String()

		if err != nil {
			//output = stderrBuffer.String()
			return
		}
	}

	output = strings.TrimRight(output, "\n")
	return
}

func (c *Client) Exists(filename string) bool {
	if _, err := c.SftpClient.Stat(filename); err == nil {
		return true
	} else {
		return false
	}
}

func (c *Client) WriteToFile(filename string, data []byte, force bool, mod os.FileMode) error {

	if err := c.connect(); err != nil {
		return err
	}

	if c.Exists(filename) && force {
		err := c.SftpClient.Remove(filename)
		if err != nil {
			return err
		}
	}

	f, err := c.SftpClient.Create(filename)
	if err != nil {
		return err
	}

	err = c.SftpClient.Chmod(filename, mod)
	if err != nil {
		return err
	}

	defer f.Close()
	if _, err = f.Write(data); err != nil {
		return err
	}
	return nil
}

func (c *Client) RenderTemplate(templateFile, destFile string, data interface{}, force bool, mode os.FileMode) error {

	tmplContent, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return err
	}

	tmpl, err := template.New(templateFile).Funcs(TemplateFuncMap).Parse(string(tmplContent))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, data); err != nil {
		return err
	}

	err = c.WriteToFile(destFile, buf.Bytes(), force, mode)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) CopyFile(localFilePath, remoteFilePath string, force bool, mode os.FileMode) error {

	if err := c.connect(); err != nil {
		return err
	}

	src, err := os.Open(localFilePath)
	if err != nil {
		return err
	}

	defer src.Close()

	_, err = src.Stat()
	if err != nil {
		return err
	}

	if c.Exists(remoteFilePath) && force {
		err = c.SftpClient.Remove(remoteFilePath)
		if err != nil {
			return err
		}
	}

	dirPath := path.Dir(remoteFilePath)
	if !c.Exists(dirPath) {
		err := c.SftpClient.MkdirAll(dirPath)
		if err != nil {
			return err
		}
	}

	dst, err := c.SftpClient.Create(remoteFilePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	err = c.SftpClient.Chmod(remoteFilePath, mode)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) readRemoteFile(path string) (*sftp.File, error) {

	if !c.Exists(path) {
		return nil, errors.New("file not found")
	}

	dst, err := c.SftpClient.Open(path)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

func (c *Client) FetchFile(remoteFilePath string, localPath string, force bool) error {

	if err := c.connect(); err != nil {
		return err
	}

	dst, err := c.readRemoteFile(remoteFilePath)
	if err != nil {
		return err
	}

	stat, err := dst.Stat()
	if err != nil {
		return err
	}

	localPath = filepath.Join(localPath, stat.Name())
	dir := filepath.Dir(localPath)
	stat, err = os.Stat(dir)
	if err != nil {
		err = os.Mkdir(dir, 0754)
		if err != nil {
			return err
		}
	}

	if force {
		_, err = os.Stat(localPath)
		if err == nil {
			err := os.Remove(localPath)
			if err != nil {
				return err
			}
		}
	}

	src, err := os.OpenFile(localPath, os.O_CREATE|os.O_WRONLY|os.O_RDONLY, 0754)
	if err != nil {
		return err
	}
	defer src.Close()

	_, err = io.Copy(src, dst)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (c *Client) MkDirs(perm os.FileMode, dirs ...string) {
	for _, dir := range dirs {
		exists, _ := c.pathExists(dir)
		if !exists {
			err := os.MkdirAll(dir, perm)
			if err != nil {
				return
			}
		}
	}
}

func (c *Client) Unarchive(localFilePath string, remoteDir string, mode os.FileMode, force bool) error {

	if err := c.connect(); err != nil {
		return err
	}

	session, err := c.SshClient.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	tmpDir := fmt.Sprintf("/tmp/%v", uuid.NewUUID())
	c.MkDirs(os.ModePerm, tmpDir)
	//defer func() {
	//	_ = os.RemoveAll(tmpDir)
	//}()

	fileList, err2 := c.DeCompress(localFilePath, tmpDir)
	if err2 != nil {
		return err2
	}

	for _, filename := range fileList {
		err = c.CopyFile(path.Join(tmpDir, filename), path.Join(remoteDir, filename), force, mode)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) Systemd(name string, state string, enable bool) error {

	if err := c.connect(); err != nil {
		return err
	}

	sysctlCmd := "systemctl %v %v"
	switch state {
	case "stop":
		sysctlCmd = fmt.Sprintf(sysctlCmd, "stop", name)
	case "start":
		sysctlCmd = fmt.Sprintf(sysctlCmd, "start", name)
	case "restart":
		sysctlCmd = fmt.Sprintf(sysctlCmd, "restart", name)
	default:
		return errors.New("state is not supported")
	}

	_, err := c.Command(sysctlCmd)
	if err != nil {
		return err
	}

	if enable {
		_, err = c.Command(fmt.Sprintf("systemctl enable %v", name))
		if err != nil {
			return err
		}
	} else {
		_, err = c.Command(fmt.Sprintf("systemctl disable %v", name))
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) Yum(name string) error {

	if err := c.connect(); err != nil {
		return err
	}
	_, err := c.Command(fmt.Sprintf("yum -y install %v", name))
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Replace(path, pattern, replace string) error {
	if err := c.connect(); err != nil {
		return err
	}

	file, err := c.readRemoteFile(path)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return errors.New("the path is dir")
	}

	if stat.Size() > TEXT_MAX_FILE_SIZE {
		return errors.New("the file is exceed the limit")
	}

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	result := re.ReplaceAllString(string(fileContent), replace)
	err = c.WriteToFile(path, []byte(result), true, stat.Mode())
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveFile(filename string) error {
	if c.Exists(filename) {
		err := c.SftpClient.Remove(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) LineInFile(path string, pattern string, insert string, line string, state string) error {
	if err := c.connect(); err != nil {
		return err
	}

	file, err := c.readRemoteFile(path)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}
	if stat.IsDir() {
		return errors.New("the path is dir")
	}

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	buff := new(bytes.Buffer)

	var text string

	if insert != "" {

		tmp := new(bytes.Buffer)
		_, err := file.WriteTo(tmp)
		if err != nil {
			return err
		}

		s := tmp.String()
		if strings.Index(s, line) != -1 {
			return nil
		}

		fileScanner = bufio.NewScanner(tmp)
		fileScanner.Split(bufio.ScanLines)

		flag := false
		for fileScanner.Scan() {
			text = fileScanner.Text()
			ok, err := regexp.MatchString(insert, text)
			if err != nil {
				fmt.Println(err)
			}
			buff.Write([]byte(text + "\n"))
			if ok && !flag {
				buff.Write([]byte(line + "\n"))
				flag = true
			}
		}

	} else {
		if pattern == "" {
			return errors.New("pattern is empty")
		}
		for fileScanner.Scan() {
			text = fileScanner.Text()
			ok, err := regexp.MatchString(pattern, text)
			if err != nil {
				fmt.Println(err)
			}
			if !ok {
				buff.Write([]byte(text + "\n"))
			}
		}

		switch state {
		case "absent":
		case "present":
			buff.Write([]byte(line))
		}
	}

	if len(buff.Bytes()) <= 0 {
		return nil
	}
	err = c.WriteToFile(path, buff.Bytes(), true, stat.Mode())
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) CreatePath(path string, fileType string, mode os.FileMode, owner string, group string) error {
	if err := c.connect(); err != nil {
		return err
	}
	if !c.Exists(path) {
		switch fileType {
		case "directory":
			err := c.SftpClient.MkdirAll(path)
			if err != nil {
				return err
			}
		case "file":
			_, err := c.SftpClient.Create(path)
			if err != nil {
				return err
			}
		}

		err := c.SftpClient.Chmod(path, mode)
		if err != nil {
			return err
		}

		_, err = c.Command(fmt.Sprintf("chown -R %v:%v %v", owner, group, path))
		if err != nil {
			return err
		}

	}

	return nil
}

func (c *Client) DeCompress(tarFile, dest string) ([]string, error) {

	var fileList []string

	srcFile, err := os.Open(tarFile)
	if err != nil {
		return nil, err
	}
	defer srcFile.Close()

	mimeType := c.getMimeType(srcFile)
	switch mimeType {
	case "application/x-gzip":
		srcFile.Seek(0, os.SEEK_SET)
		gr, err := gzip.NewReader(srcFile)
		if err != nil {
			return nil, err
		}
		defer gr.Close()
		fileList, err = c.deTarCompress(gr, dest)
	case "application/octet-stream":
		srcFile.Seek(0, os.SEEK_SET)
		fileList, err = c.deTarCompress(srcFile, dest)
	}

	return fileList, nil
}

func (c *Client) deTarCompress(srcFile io.Reader, dest string) ([]string, error) {
	fileList := make([]string, 0)
	tr := tar.NewReader(srcFile)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		filename := dest + "/" + hdr.Name
		file, err := c.createFile(filename)
		if err != nil {
			continue
		}
		io.Copy(file, tr)
		fileList = append(fileList, hdr.Name)
		file.Close()
	}

	return fileList, nil
}

func (c *Client) createFile(name string) (*os.File, error) {
	path := string([]rune(name)[0:strings.LastIndex(name, "/")])
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}

func (c *Client) getMimeType(f *os.File) string {
	buffer := make([]byte, 512)
	_, _ = f.Read(buffer)
	contentType := http.DetectContentType(buffer)
	return contentType
}

func (c *Client) Close() {
	if c.SshClient != nil {
		c.SshClient.Close()
	}
	if c.SftpClient != nil {
		c.SftpClient.Close()
	}
}

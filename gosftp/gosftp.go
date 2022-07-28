package gosftp

import (
	"bytes"
	"fmt"
	"gosail/gossh"
	"gosail/model"
	"gosail/utils"
	"io/ioutil"
	"sync"

	"os"
	"path"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func NewClient(sshClient *ssh.Client) (*sftp.Client, error) {
	// create sftp client
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, err
	}
	return sftpClient, nil
}

func uploadFile(sftpClient *sftp.Client, localFilePath string, remotePath string, buffer *bytes.Buffer) error {
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		return fmt.Errorf("os.Open error, %s", err)
	}
	defer srcFile.Close()

	var remoteFileName = utils.GetPathLastName(localFilePath)
	err = sftpClient.MkdirAll(remotePath)
	if err != nil {
		return fmt.Errorf("sftpClient.Mkdir error, %s", err)
	}

	dstFile, err := sftpClient.Create(path.Join(remotePath, remoteFileName))
	if err != nil {
		return fmt.Errorf("sftpClient.Create error, %s", err)
	}
	defer dstFile.Close()

	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		return fmt.Errorf("readAll error, %s", err)
	}
	_, err = dstFile.Write(ff)
	if err != nil {
		return fmt.Errorf("write error, %s", err)
	}
	buffer.WriteString(localFilePath)
	buffer.WriteString(" -> ")
	buffer.WriteString(path.Join(remotePath, remoteFileName))
	buffer.WriteString("\n")
	return nil
}

func uploadDirectory(sftpClient *sftp.Client, localPath string, remotePath string, buffer *bytes.Buffer) error {
	localFiles, err := ioutil.ReadDir(localPath)
	if err != nil {
		return fmt.Errorf("readDir error, %s", err)
	}
	sftpClient.Mkdir(remotePath)
	for _, backupDir := range localFiles {
		localFilePath := path.Join(localPath, backupDir.Name())
		remoteFilePath := path.Join(remotePath, backupDir.Name())

		if backupDir.IsDir() {
			err = sftpClient.Mkdir(remoteFilePath)
			if err != nil {
				return fmt.Errorf("sftpClient.Mkdir error, %s", err)
			}
			err := uploadDirectory(sftpClient, localFilePath, remoteFilePath, buffer)
			if err != nil {
				return fmt.Errorf("sftpClient.UploadDirectory error, %s", err)
			}
		} else {
			uploadFile(sftpClient, path.Join(localPath, backupDir.Name()), remotePath, buffer)
		}
	}
	buffer.WriteString(localPath)
	buffer.WriteString(" -> ")
	buffer.WriteString(remotePath)
	buffer.WriteString("\n")
	return nil
}

func upload(sftpClient *sftp.Client, localPath string, remotePath string) (string, error) {
	local, err := os.Stat(localPath)
	if err != nil {
		return "", fmt.Errorf("os.Stat error, %s", err)
	}
	var buffer bytes.Buffer
	if local.IsDir() {
		var remoteDirName = utils.GetPathLastName(localPath)
		remotePath = path.Join(remotePath, remoteDirName)
		err = uploadDirectory(sftpClient, localPath, remotePath, &buffer)
		if err != nil {
			return buffer.String(), err
		}
	} else {
		err = uploadFile(sftpClient, localPath, remotePath, &buffer)
		if err != nil {
			return buffer.String(), err
		}
	}
	return buffer.String(), err
}

func downloadFile(sftpClient *sftp.Client, remotePath string, localFilePath string, buffer *bytes.Buffer) error {
	srcFile, err := sftpClient.Open(remotePath)
	if err != nil {
		return fmt.Errorf("sftpClient.Open error, %s", err)
	}
	defer srcFile.Close()

	var localFileName = utils.GetPathLastName(remotePath)
	err = os.MkdirAll(localFilePath, 0777)
	if err != nil {
		return fmt.Errorf("os.Mkdir error, %s", err)
	}

	dstFile, err := os.Create(path.Join(localFilePath, localFileName))
	if err != nil {
		return fmt.Errorf("sftpClient.Create error, %s", err)
	}
	defer dstFile.Close()

	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		return fmt.Errorf("readAll error, %s", err)
	}
	_, err = dstFile.Write(ff)
	if err != nil {
		return fmt.Errorf("write error, %s ", err)
	}
	buffer.WriteString(remotePath)
	buffer.WriteString(" -> ")
	buffer.WriteString(path.Join(localFilePath, localFileName))
	buffer.WriteString("\n")
	return nil
}

func downloadDirectory(sftpClient *sftp.Client, remotePath string, localPath string, buffer *bytes.Buffer) error {
	remoteFiles, err := sftpClient.ReadDir(remotePath)
	if err != nil {
		return fmt.Errorf("sftpClient.ReadDir error, %s", err)
	}

	err = os.MkdirAll(localPath, 0777)
	if err != nil {
		return fmt.Errorf("os.Mkdir error, %s", err)
	}
	for _, backupDir := range remoteFiles {
		localFilePath := path.Join(localPath, backupDir.Name())
		remoteFilePath := path.Join(remotePath, backupDir.Name())
		if backupDir.IsDir() {
			sftpClient.Mkdir(localFilePath)
			downloadDirectory(sftpClient, remoteFilePath, localFilePath, buffer)
		} else {
			downloadFile(sftpClient, path.Join(remotePath, backupDir.Name()), localPath, buffer)
		}
	}
	buffer.WriteString(remotePath)
	buffer.WriteString(" -> ")
	buffer.WriteString(localPath)
	buffer.WriteString("\n")
	return nil
}

func download(sftpClient *sftp.Client, remotePath string, localPath string) (string, error) {
	remote, err := sftpClient.Stat(remotePath)
	if err != nil {
		return "", fmt.Errorf("sftpClient.Stat error, '%s' %s", remotePath, err)
	}
	var buffer bytes.Buffer
	if remote.IsDir() {
		var localDirName = utils.GetPathLastName(remotePath)
		localPath = path.Join(localPath, localDirName)
		err = downloadDirectory(sftpClient, remotePath, localPath, &buffer)
		if err != nil {
			return buffer.String(), err
		}
	} else {
		err = downloadFile(sftpClient, remotePath, localPath, &buffer)
		if err != nil {
			return buffer.String(), err
		}
	}
	return buffer.String(), nil
}

func ClientUpload(chLimit chan struct{}, host model.SSHHost, clientConfig *model.ClientConfig, srcPath, destPath string, ch chan model.RunResult, wg *sync.WaitGroup) {
	defer wg.Done()
	client, err := gossh.Connect(host.Username, host.Password, host.Host, host.Key, host.Port, clientConfig.CipherList, clientConfig.KeyExchangeList)

	var sftpResult model.RunResult
	sftpResult.Host = host.Host
	sftpResult.Username = host.Username

	if err != nil {
		sftpResult.Success = false
		sftpResult.Result = fmt.Sprintf("<%s>\n", err.Error())
		ch <- sftpResult
		return
	}

	sftpClient, err := NewClient(client)
	if err != nil {
		sftpResult.Success = false
		sftpResult.Result = fmt.Sprintf("<%s>\n", err.Error())
		ch <- sftpResult
		return
	}
	defer sftpClient.Close()

	result, err := upload(sftpClient, srcPath, destPath)
	if err != nil {
		sftpResult.Success = false
		sftpResult.Result = fmt.Sprintf("%s\n", err.Error())
		ch <- sftpResult
		return
	}
	sftpResult.Success = true
	sftpResult.Result = result
	ch <- sftpResult
}

func ClientDownload(chLimit chan struct{}, host model.SSHHost, clientConfig *model.ClientConfig, srcPath, destPath string, ch chan model.RunResult, wg *sync.WaitGroup) {
	defer wg.Done()
	client, err := gossh.Connect(host.Username, host.Password, host.Host, host.Key, host.Port, clientConfig.CipherList, clientConfig.KeyExchangeList)

	var sftpResult model.RunResult
	sftpResult.Host = host.Host
	sftpResult.Username = host.Username

	if err != nil {
		sftpResult.Success = false
		sftpResult.Result = fmt.Sprintf("<%s>\n", err.Error())
		ch <- sftpResult
		return
	}

	sftpClient, err := NewClient(client)
	if err != nil {
		sftpResult.Success = false
		sftpResult.Result = fmt.Sprintf("<%s>\n", err.Error())
		ch <- sftpResult
		return
	}
	defer sftpClient.Close()

	result, err := download(sftpClient, srcPath, destPath)
	if err != nil {
		sftpResult.Success = false
		sftpResult.Result = fmt.Sprintf("%s\n", err.Error())
		ch <- sftpResult
		return
	}
	sftpResult.Success = true
	sftpResult.Result = result
	ch <- sftpResult
}

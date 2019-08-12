package utils

import (
	"archive/zip"
	"io"
	"os"
)

/**
@files：需要压缩的文件
@compreFile：压缩之后的文件
*/
func ZipCompress(files []*os.File, compreFile *os.File) (err error) {
	zw := zip.NewWriter(compreFile)
	defer zw.Close()
	for _, file := range files {
		err := compress_zip(file, zw)
		if err != nil {
			return err
		}
		file.Close()
	}
	return nil
}

/**
功能：压缩文件
@file:压缩文件
@prefix：压缩文件内部的路径
@tw：写入压缩文件的流
*/
func compress_zip(file *os.File, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	// 获取压缩头信息
	head, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	// 指定文件压缩方式 默认为 Store 方式 该方式不压缩文件 只是转换为zip保存
	head.Method = zip.Deflate
	fw, err := zw.CreateHeader(head)
	if err != nil {
		return err
	}
	// 写入文件到压缩包中
	_, err = io.Copy(fw, file)
	file.Close()
	if err != nil {
		return err
	}
	return nil
}

/**
@tarFile：压缩文件路径
@dest：解压文件夹
*/
func ZipDeCompressByPath(tarFile, dest string) error {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	return ZipDeCompress(srcFile, dest)
}

/**
@zipFile：压缩文件
@dest：解压之后文件保存路径
*/
func ZipDeCompress(srcFile *os.File, dest string) error {
	zipFile, err := zip.OpenReader(srcFile.Name())
	if err != nil {
		return err
	}
	defer zipFile.Close()
	for _, innerFile := range zipFile.File {
		info := innerFile.FileInfo()
		if info.IsDir() {
			err = os.MkdirAll(innerFile.Name, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}
		srcFile, err := innerFile.Open()
		if err != nil {
			continue
		}
		defer srcFile.Close()
		newFile, err := os.Create(innerFile.Name)
		if err != nil {
			continue
		}
		io.Copy(newFile, srcFile)
		newFile.Close()
	}
	return nil
}

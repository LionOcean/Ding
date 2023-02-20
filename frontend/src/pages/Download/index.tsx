import { Button, Space, Table, Input, message } from 'antd';
import { CloudUploadOutlined, DeleteOutlined } from '@ant-design/icons';
import React, { useCallback, useEffect, useRef, useState } from "react";
import { ColumnsType } from 'antd/es/table';
import { DownloadFile, LocalIPAddr, SaveFileDialog } from '@wailsjs/go/main/App';
import { isEqualLAN } from '@utils/index';
import { decrypt, encrypt } from '@utils/crypto';
import { OpenDirDialog } from "../../../wailsjs/go/main/App";
const { Search } = Input;
interface DataType {
  key: React.Key;
  file: string;
  path: string;
  size: number;
  name: string;
}

const columns: ColumnsType<DataType> = [
  {
    title: '文件名',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '文件路径',
    dataIndex: 'path',
    key: 'path',
  },
  {
    title: '大小/kb',
    dataIndex: 'size',
    key: 'size',
  },
];

export default function Download() {
  const [files, setFiles] = useState<Array<DataType>>([]);
  const [selectedFiles, setSelectedFiled] = useState<Array<DataType>>([]);
  const [localIp, setLocalIp] = useState<string[]>([]);
  const [remoteIp, setRemoteIp] = useState<string[]>([]);
  const downloadLoading = useRef(false)
  const onDownloadFiles = useCallback(async () => {
    try {
      if(downloadLoading.current) return
      downloadLoading.current = true
      let localUrl = await OpenDirDialog({
        Title: "选择下载目录",
      })
      if(!localUrl) {
        downloadLoading.current = false
        return message.error("未选择要下载的目录")
      }
      for (const file of selectedFiles) {
        const prefix = file.name.slice(file.name.lastIndexOf("."))
        const remoteUrl = `http://${remoteIp.join(':')}/download?path=${file.path}`
        let res = await DownloadFile(remoteUrl, localUrl + "/" + file.name);
        if(!res!) {
          message.success(`文件${file.name}下载成功`)
        } else {
          message.error(`文件${file.name}下载失败`)
        }
        downloadLoading.current = false
      }
    } catch (e) {
      downloadLoading.current = false
      console.log(e)
    }
  }, []);
  useEffect(() => {
    (async () => {
      let [ip, port] = await LocalIPAddr();
      setLocalIp([ip, port]);
    })();
    return;
  }, []);
  const onDeleteFiles = useCallback(() => {}, []);
  const onSearch = (value: string) => {
    if (!value) {
      return message.error('请传入hash');
    }
    const remoteAddr = decrypt(value).split(',');
    setRemoteIp(remoteAddr);
    console.log(remoteAddr, localIp, remoteIp);
    if (!isEqualLAN(localIp[0], remoteAddr[0])) {
      return message.error('发送端IP与本机IP不属于同一局域网');
    }
    const requestUrl = `http://${remoteAddr.join(':')}/list`;

    fetch(requestUrl).then(async (res) => {
      const result = await res.json();
      if (result.code === 200) {
        setFiles(result.data);
      } else {
        message.error(result.data);
      }
    });
  };
  return (
    <div className='download'>
      <Search
        placeholder='请输入传输密钥'
        onSearch={onSearch}
        style={{
          width: 'calc(100% - 48px)',
          margin: '16px 0',
        }}
      ></Search>
      <Space>
        <Button icon={<CloudUploadOutlined />} size={'middle'} onClick={onDownloadFiles}>
          下载
        </Button>
        <Button icon={<DeleteOutlined />} size={'middle'} onClick={onDeleteFiles}>
          移除
        </Button>
      </Space>
      <Table
        rowSelection={{
          type: 'checkbox',
          onChange: (selectedRowKeys: React.Key[], selectedRows: DataType[]) => {
            setSelectedFiled(selectedRows);
          },
          getCheckboxProps: (record: DataType) => ({
            disabled: false,
            ...record,
          }),
        }}
        columns={columns}
        dataSource={files}
        pagination={false}
        scroll={{ y: 500 }}
        style={{
          marginTop: 15,
        }}
      />
    </div>
  );
}

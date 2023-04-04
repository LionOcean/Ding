import { Button, Space, Table, Input, message, Progress } from 'antd';
import { FileAddOutlined } from '@ant-design/icons';
import React, { useCallback, useEffect, useRef, useState } from 'react';
import { ColumnsType } from 'antd/es/table';
import { OpenDirDialog } from '@wailsjs/go/main/App';
import { DownloadFile, ReceivingFiles } from '@wailsjs/go/transfer/ReceivePeer';
import { LocalIPAddr } from '@wailsjs/go/transfer/Peer';
import { EventsOn, EventsOff } from '@wailsjs/runtime/runtime';
import { isEqualLAN, calcByteUnit } from '@/utils';
import { decrypt } from '@utils/crypto';

import './index.less';

const { Search } = Input;

interface DataType {
  key: React.Key;
  file: string;
  path: string;
  size: number;
  name: string;
  sizeUnit: string;
  // 下载进度相关
  download: {
    progress: number;
    status: 'success' | 'normal' | 'exception' | 'active' | undefined;
  };
}

const columns: ColumnsType<DataType> = [
  {
    title: '文件名',
    dataIndex: 'name',
    align: 'center',
    ellipsis: {
      showTitle: true,
    },
  },
  // 下载先不让用户知道发送者电脑的目录信息
  // {
  //   title: '文件路径',
  //   dataIndex: 'path',
  // },
  {
    title: '进度',
    render(_, record) {
      return <Progress status={record.download.status} percent={record.download.progress} steps={10} size='small' />;
    },
  },
  {
    title: '大小',
    dataIndex: 'sizeUnit',
    align: 'center',
  },
];

export default function Download() {
  const [files, setFiles] = useState<Array<DataType>>([]);
  const [selectedFiles, setSelectedFiled] = useState<Array<DataType>>([]);
  const [localIp, setLocalIp] = useState<string[]>([]);
  const [remoteIp, setRemoteIp] = useState<string[]>([]);
  const downloadLoading = useRef(false);

  // 监听下载进度
  EventsOff('EVENT_DOWN_PROGRESS');
  EventsOn('EVENT_DOWN_PROGRESS', (data) => {
    const { name, finished, total } = data;
    const downloadFiles = files.map((file) => {
      if (file.name === name) {
        file.download.progress = Math.floor(finished / total) * 100;
      }
      return file;
    });
    setFiles(downloadFiles);
    console.log('文件下载进度: ', data);
  });

  const onDownloadFiles = useCallback(async () => {
    try {
      if (downloadLoading.current) return;
      downloadLoading.current = true;
      let localUrl = await OpenDirDialog({
        Title: '选择下载目录',
      });
      if (!localUrl) {
        downloadLoading.current = false;
        return message.error('未选择要下载的目录');
      }
      console.log(selectedFiles);

      // 每次下载前清空所有文件的下载进度
      const downloadFiles = selectedFiles.map((file) => {
        file.download.progress = 0;
        return file;
      });
      setSelectedFiled(downloadFiles);

      for (const file of selectedFiles) {
        const remoteUrl = remoteIp.join(':');
        DownloadFile(remoteUrl, file, localUrl + '/' + file.name)
          .then(() => {
            // message.success(`文件${file.name}下载成功`);
          })
          .catch((e) => {
            const downloadFiles = selectedFiles.map((dFile) => {
              if (dFile.name === file.name) dFile.download.status = 'exception';
              return dFile;
            });
            setSelectedFiled(downloadFiles);
            message.error(`文件${file.name}下载失败 ${e}`);
          });
      }
      downloadLoading.current = false;
    } catch (e) {
      downloadLoading.current = false;
      console.log(e);
    }
  }, [selectedFiles]);

  useEffect(() => {
    (async () => {
      let [ip, port] = await LocalIPAddr();
      setLocalIp([ip, port]);
    })();
    return;
  }, []);

  const onSearch = (value: string) => {
    if (!value) {
      return message.error('请传入hash');
    }
    const remoteAddr = decrypt(value).split(',');
    setRemoteIp(remoteAddr);
    // if (!isEqualLAN(localIp[0], remoteAddr[0])) {
    //   return message.error('发送端IP与本机IP不属于同一局域网');
    // }
    ReceivingFiles(remoteAddr.join(':'))
      .then((res) => {
        const result = JSON.parse(res);
        if (result.code === 200) {
          setFiles(
            result.data.map((ele: DataType) => {
              const r = calcByteUnit(ele.size);
              return {
                ...ele,
                sizeUnit: r.join(''),
                key: ele.name,
                download: {
                  progress: 0,
                  status: 'normal',
                },
              };
            })
          );
        } else {
          message.error(result.data);
        }
      })
      .catch((e) => {
        console.log(e);
        message.error('连接超时, 请重试');
      });
  };
  return (
    <div className='download'>
      <Space>
        <Search placeholder='请输入传输密钥' onSearch={onSearch}></Search>
        <Button icon={<FileAddOutlined />} size={'middle'} onClick={onDownloadFiles}>
          下载
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

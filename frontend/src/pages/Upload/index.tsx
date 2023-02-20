import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { UploadFiles, RemoveFiles, LogTransferFiles, LocalIPAddr, StartP2PServer } from '@wailsjs/go/main/App';
import { transfer } from '@wailsjs/go/models';
import { Button, Input, Modal, Space, Table, message } from 'antd';
import { CloudUploadOutlined, CopyOutlined, DeleteOutlined } from '@ant-design/icons';
import { ColumnsType } from 'antd/es/table';
import ClipboardJS from 'clipboard';
import { encrypt } from '@utils/crypto';
import './index.less';

interface DataType {
  key: React.Key;
  name: string;
  path: string;
  size: number;
}

const columns: ColumnsType<DataType> = [
  {
    title: '文件名',
    dataIndex: 'name',
  },
  {
    title: '文件路径',
    render: (_, record) => {
      return (
        <a
          color='geekblue'
          onClick={() => {
            Modal.info({
              content: record.path,
              icon: null,
              centered: true,
            });
          }}
        >
          查看
        </a>
      );
    },
  },
  {
    title: '大小/kb',
    dataIndex: 'size',
  },
];

// 包装files list数组
function wrapFiles(raw: transfer.TransferFile[]): DataType[] {
  const target: DataType[] = raw.map((item, index) => {
    const { size } = item;
    return {
      key: index,
      ...item,
      size: Math.ceil(size / 1024),
    };
  });
  return target;
}

export default function Upload() {
  const [files, setFiles] = useState<Array<DataType>>([]);
  const [selectedFiles, setSelectedFiled] = useState<Array<DataType>>([]);
  const [ipAddr, setIpAddr] = useState('');

  const uploading = useRef(false);

  const [messageApi, contextHolder] = message.useMessage();

  const selectedRowKeys = useMemo(() => {
    return selectedFiles.map((item) => item.key);
  }, [selectedFiles]);

  const filterFiles = (targetFiles: DataType[]) => {
    if (!targetFiles?.length) {
      return null;
    }
    return files.filter((file) => targetFiles.findIndex((item) => item.name === file.name) < 0);
  };

  // 每次先加载列表
  LogTransferFiles().then((res: transfer.TransferFile[]) => {
    setFiles(wrapFiles(res));
  })

  const onUploadFiles = useCallback(async () => {
    if (uploading.current) return;
    uploading.current = true;
    try {
      const res = await UploadFiles({});
      setFiles(wrapFiles(res));
    } catch (error) {}
    uploading.current = false;
  }, []);

  const onDeleteFiles = useCallback(async () => {
    if (!selectedFiles.length) {
      messageApi.open({
        type: 'warning',
        content: '请选中删除的文件',
      });
      return;
    }
    try {
      const res = await RemoveFiles(selectedFiles);
      !!res && setFiles(wrapFiles(res));
      setSelectedFiled([]);
    } catch (error) {
      console.error(error);
      messageApi.open({
        type: 'error',
        content: '删除失败,请重试',
      });
    }
  }, [selectedFiles]);

  // 初始化 clipboard
  useEffect(() => {
    const clipboard = new ClipboardJS('.clipboardBtn');

    clipboard.on('success', (e) => {
      messageApi.open({
        type: 'success',
        content: '复制成功',
      });
    });

    clipboard.on('error', (e) => {});
  }, []);

  useEffect(() => {
    try {
      (async () => {
        const addr = await LocalIPAddr();
        const targetAddr = encrypt(addr.join(','));
        setIpAddr(targetAddr);
      })();
    } catch (error) {
      console.log(error);
    }
  }, []);

  useEffect(() => {
    try {
      StartP2PServer().catch(err => {
        console.log(err);
      });
    } catch (error) {
      console.log(error);
    }
  }, []);

  return (
    <>
      {contextHolder}
      <div className='upload'>
        <Space style={{ width: '100%', paddingLeft: '8px' }}>
          <Input value={ipAddr} style={{ width: '250px', backgroundColor: '#fff', color: '#666' }} readOnly />
          <Button className='clipboardBtn' data-clipboard-text={ipAddr}>
            <CopyOutlined style={{ fontSize: '18px' }} />
          </Button>
        </Space>
        <Space style={{ width: '100%', paddingLeft: '8px', marginTop: '15px' }}>
          <Button icon={<CloudUploadOutlined />} size={'middle'} onClick={onUploadFiles}>
            上传
          </Button>
          <Button icon={<DeleteOutlined />} size={'middle'} onClick={onDeleteFiles}>
            移除
          </Button>
        </Space>
        <Table
          rowSelection={{
            type: 'checkbox',
            onChange: (_, selectedRows: DataType[]) => {
              setSelectedFiled(selectedRows);
            },
            selectedRowKeys,
          }}
          columns={columns}
          dataSource={files}
          pagination={false}
          scroll={{ y: 500 }}
          style={{
            marginTop: 15,
          }}
        ></Table>
      </div>
    </>
  );
}

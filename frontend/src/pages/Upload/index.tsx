import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { UploadFiles, Remove, List, StartServer } from '@wailsjs/go/transfer/SendPeer';
import { LocalIPAddr } from '@wailsjs/go/transfer/Peer';
import { transfer } from '@wailsjs/go/models';
import { Button, Input, Space, Table, message, Tooltip } from 'antd';
import { FolderOpenOutlined, CopyOutlined, DeleteOutlined } from '@ant-design/icons';
import { ColumnsType } from 'antd/es/table';
import ClipboardJS from 'clipboard';
import { encrypt } from '@utils/crypto';
import { calcByteUnit } from '@utils/index';

import './index.less';

interface DataType {
  key: React.Key;
  name: string;
  path: string;
  size: number;
  sizeUnit: string;
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
  {
    title: '文件路径',
    align: 'center',
    render: (_, record) => {
      return (
        <Tooltip title={record.path} trigger='click' color='#108ee9'>
          <a color='geekblue'>查看</a>
        </Tooltip>
      );
    },
  },
  {
    title: '大小',
    dataIndex: 'sizeUnit',
    align: 'center',
  },
];

// 包装files list数组
function wrapFiles(raw: transfer.TransferFile[]): DataType[] {
  const target: DataType[] = raw.map((item, index) => {
    const { size } = item;
    const r = calcByteUnit(size);
    return {
      key: index,
      ...item,
      sizeUnit: r.join(''),
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

  // 每次先加载列表
  List().then((res: transfer.TransferFile[]) => {
    setFiles(wrapFiles(res));
  });

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
      const res = await Remove(selectedFiles);
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
      StartServer().catch((err) => {
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
        <Space style={{ width: '100%', paddingLeft: '8px', justifyContent: 'center' }}>
          <Input value={ipAddr} style={{ width: '250px', backgroundColor: '#fff', color: '#666' }} readOnly />
          <Button className='clipboardBtn' data-clipboard-text={ipAddr}>
            <CopyOutlined style={{ fontSize: '18px' }} />
          </Button>
        </Space>
        <Space style={{ width: '100%', paddingLeft: '8px', marginTop: '15px', justifyContent: 'center' }}>
          <Button icon={<FolderOpenOutlined />} size={'middle'} onClick={onUploadFiles}>
            文件
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

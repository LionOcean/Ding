import { useCallback, useRef, useState } from 'react';
import { LogFiles, UploadFiles } from '../../../wailsjs/go/main/App';
import { Button, Space, Table } from 'antd';
import { CloudUploadOutlined, DeleteOutlined } from '@ant-design/icons';
import { ColumnsType } from 'antd/es/table';

import './index.less';

interface DataType {
  key: React.Key;
  file: string;
  path: string;
  size: number;
}

const columns: ColumnsType<DataType> = [
  {
    title: '文件名',
    dataIndex: 'file',
  },
  {
    title: '文件路径',
    dataIndex: 'path',
  },
  {
    title: '大小/kb',
    dataIndex: 'size',
  },
];

export default function Upload() {
  const [files, setFiles] = useState<Array<DataType>>([]);
  const [selectedFiles, setSelectedFiled] = useState<Array<DataType>>([]);
  const uploading = useRef(false);

  const onUploadFiles = useCallback(async () => {
    if (uploading.current) return;
    uploading.current = true;
    try {
      await UploadFiles({});
      const res = await LogFiles();
      const target: DataType[] = res.map((item, index) => {
        const { path, size } = item;
        return {
          key: index,
          file: path.split('\\').pop() || path,
          path,
          size: Math.ceil(size / 1024),
        };
      });
      setFiles(target);
      console.log('上传文件成功.', res);
    } catch (error) {}
    uploading.current = false;
  }, []);

  const onDeleteFiles = useCallback(async () => {}, [selectedFiles]);

  return (
    <div className='upload'>
      <Space>
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

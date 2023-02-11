import {Button, Space, Table, Input} from "antd";
import {CloudUploadOutlined, DeleteOutlined} from "@ant-design/icons";
import {useCallback, useState} from "react";
import {ServerIPAddr} from "../../../wailsjs/go/main/App";
import {ColumnsType} from "antd/es/table";
const {Search} = Input
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

export default function Download() {
    const [files, setFiles] = useState<Array<DataType>>([]);
    const [selectedFiles, setSelectedFiled] = useState<Array<DataType>>([]);
    const onDownloadFiles = useCallback(async () => {
        let ip= await ServerIPAddr()
        console.log(ip)
    }, []);
    const onDeleteFiles = useCallback(() => {

    }, [])
    const onSearch = (value: string) => {
        console.log(value)
    }
  return (
      <div className='download'>
              <Search placeholder="请输入传输密钥" onSearch={onSearch} style={{
                  width: 'calc(100% - 48px)',
                  margin: '16px 0'
              }}></Search>
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

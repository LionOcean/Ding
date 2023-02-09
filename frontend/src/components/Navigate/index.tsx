import { useNavigate } from 'react-router-dom';
import { DownloadOutlined, UploadOutlined } from '@ant-design/icons';
import './index.less';
import { useCallback } from 'react';

export default function Navigate() {
  const navigate = useNavigate();

  const goto = useCallback((path: string) => {
    navigate(path);
  }, []);

  return (
    <div className='navigate'>
      <div className='navigate__item' role='button' tabIndex={-1} onClick={() => goto('/upload')}>
        <UploadOutlined style={{ fontSize: '36px' }} />
        <span className='navigate__text'>上传</span>
      </div>
      <div className='navigate__item' role='button' tabIndex={-1} onClick={() => goto('/download')}>
        <DownloadOutlined style={{ fontSize: '36px' }} />
        <span className='navigate__text'>下载</span>
      </div>
    </div>
  );
}

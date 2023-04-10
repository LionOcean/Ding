import { useNavigate, useLocation } from 'react-router-dom';
import { ShareAltOutlined, SendOutlined } from '@ant-design/icons';
import './index.less';
import { useCallback } from 'react';

export default function Navigate() {
  const navigate = useNavigate();

  const goto = useCallback((path: string) => {
    navigate(path);
  }, []);

  const { pathname } = useLocation()
  // 非home页不展示上传/下载跳转按钮
  const show = pathname === '/home';

  if (show)  {
    return (
      <div className='navigate'>
        <div className='navigate__item' role='button' tabIndex={-1} onClick={() => goto('/upload')}>
          <SendOutlined style={{ fontSize: '28px' }} />
          <span className='navigate__text'>发送</span>
        </div>
        <div className='navigate__item' role='button' tabIndex={-1} onClick={() => goto('/download')}>
          <ShareAltOutlined style={{ fontSize: '30px' }} />
          <span className='navigate__text'>接收</span>
        </div>
      </div>
    );
  } else {
    return <></>
  }
}

import { createRoot } from 'react-dom/client';

import App from './App';

import './style.less';

const container = document.getElementById('root');

const root = createRoot(container!);

root.render(<App />);

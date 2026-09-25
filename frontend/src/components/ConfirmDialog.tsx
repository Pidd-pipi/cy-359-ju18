// 共享组件：确认对话框封装
import { Modal } from 'antd'
import { ExclamationCircleOutlined } from '@ant-design/icons'

export function confirmAction(title: string, content: string, onOk: () => void | Promise<void>) {
  Modal.confirm({
    title,
    icon: <ExclamationCircleOutlined />,
    content,
    okText: '确认',
    cancelText: '取消',
    onOk,
  })
}

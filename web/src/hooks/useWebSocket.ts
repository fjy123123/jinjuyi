import { useCallback, useEffect, useRef, useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { addMessage, updateMessageReadStatus } from '@/store/chat'
import { RootState } from '@/store'

interface WebSocketHook {
  isConnected: boolean
  sendMessage: (data: any) => void
}

export const useWebSocket = (): WebSocketHook => {
  const wsRef = useRef<WebSocket | null>(null)
  const [isConnected, setIsConnected] = useState(false)
  const dispatch = useDispatch()
  const { user } = useSelector((state: RootState) => state.auth)

  const WS_BASE_URL = import.meta.env.VITE_WS_BASE_URL || 'ws://localhost:8080'

  const connect = useCallback(() => {
    try {
      const token = localStorage.getItem('token')
      const url = `${WS_BASE_URL}/ws?token=${token}`
      
      wsRef.current = new WebSocket(url)

      wsRef.current.onopen = () => {
        console.log('WebSocket connected')
        setIsConnected(true)
      }

      wsRef.current.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          
          // 处理新消息
          if (data.type === 'new_message') {
            const message = data.data
            dispatch(addMessage({
              type: message.group_id ? 2 : 1,
              targetId: message.group_id || message.receiver_id,
              message
            }))
          }
          
          // 处理消息撤回
          if (data.type === 'recall_message') {
            console.log('Recall message:', data.data)
          }
          
          // 处理已读回执
          if (data.type === 'read_receipt') {
            const { user_id, target_id, type, message_ids } = data.data
            if (user_id !== user?.id) {
              dispatch(updateMessageReadStatus({
                convType: type,
                targetId: target_id,
                messageIds: message_ids,
                readerId: user_id
              }))
            }
          }
        } catch (error) {
          console.error('WebSocket message parse error:', error)
        }
      }

      wsRef.current.onclose = () => {
        console.log('WebSocket disconnected')
        setIsConnected(false)
        // 自动重连
        setTimeout(() => {
          connect()
        }, 3000)
      }

      wsRef.current.onerror = (error) => {
        console.error('WebSocket error:', error)
        setIsConnected(false)
      }

    } catch (error) {
      console.error('WebSocket connection error:', error)
    }
  }, [WS_BASE_URL, dispatch, user?.id])

  const sendMessage = useCallback((data: any) => {
    if (wsRef.current && isConnected) {
      wsRef.current.send(JSON.stringify(data))
    } else {
      console.warn('WebSocket not connected')
    }
  }, [isConnected])

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (token) {
      connect()
    }

    return () => {
      if (wsRef.current) {
        wsRef.current.close()
      }
    }
  }, [connect])

  return {
    isConnected,
    sendMessage
  }
}
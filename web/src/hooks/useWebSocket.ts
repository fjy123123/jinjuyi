import { useCallback, useEffect, useRef, useState } from 'react'
import { useDispatch } from 'react-redux'
import { addMessage } from '@/store/chat'

interface WebSocketHook {
  isConnected: boolean
  sendMessage: (data: any) => void
}

export const useWebSocket = (): WebSocketHook => {
  const wsRef = useRef<WebSocket | null>(null)
  const [isConnected, setIsConnected] = useState(false)
  const dispatch = useDispatch()

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
          
          // 处理不同类型的消息
          if (data.type === 'message') {
            const message = data.payload
            dispatch(addMessage({
              type: message.group_id ? 2 : 1,
              targetId: message.group_id || message.sender_id,
              message
            }))
          }
          
          if (data.type === 'redpacket') {
            // 处理红包消息
            console.log('Red packet message:', data.payload)
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
  }, [WS_BASE_URL, dispatch])

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
import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit'
import { api } from '@/services/api'

interface User {
  id: number
  username: string
  nickname: string
  avatar?: string
}

interface Message {
  id: number
  sender_id: number
  content: string
  message_type: number
  created_at: string
  is_recall?: boolean
  sender?: User
}

interface Conversation {
  target_id: number
  target_name: string
  target_avatar?: string
  type: 1 | 2
  last_message?: Message
  last_message_at?: string
  unread_count: number
}

interface ChatState {
  conversations: Conversation[]
  currentConversation: Conversation | null
  messages: Record<string, Message[]>
  loading: boolean
}

const initialState: ChatState = {
  conversations: [],
  currentConversation: null,
  messages: {},
  loading: false
}

export const fetchConversations = createAsyncThunk(
  'chat/fetchConversations',
  async (_, { rejectWithValue }) => {
    try {
      const res = await api.conversation.getList()
      if (res.code === 0) {
        return res.data
      }
      throw new Error(res.message)
    } catch (error: any) {
      return rejectWithValue(error.message)
    }
  }
)

export const fetchMessages = createAsyncThunk(
  'chat/fetchMessages',
  async ({ type, targetId }: { type: 1 | 2; targetId: number }, { rejectWithValue }) => {
    try {
      let res
      if (type === 1) {
        res = await api.message.getPrivate(targetId)
      } else {
        res = await api.message.getGroup(targetId)
      }
      if (res.code === 0) {
        return { type, targetId, messages: res.data }
      }
      throw new Error(res.message)
    } catch (error: any) {
      return rejectWithValue(error.message)
    }
  }
)

const chatSlice = createSlice({
  name: 'chat',
  initialState,
  reducers: {
    setCurrentConversation: (state, action: PayloadAction<Conversation>) => {
      state.currentConversation = action.payload
    },
    addMessage: (state, action: PayloadAction<{ type: 1 | 2; targetId: number; message: Message }>) => {
      const key = action.payload.type === 1 
        ? `private_${action.payload.targetId}` 
        : `group_${action.payload.targetId}`
      if (state.messages[key]) {
        state.messages[key].push(action.payload.message)
      } else {
        state.messages[key] = [action.payload.message]
      }
    },
    setConversations: (state, action: PayloadAction<Conversation[]>) => {
      state.conversations = action.payload
    },
    setMessages: (state, action: PayloadAction<{ key: string; messages: Message[] }>) => {
      state.messages[action.payload.key] = action.payload.messages
    }
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchConversations.pending, (state) => {
        state.loading = true
      })
      .addCase(fetchConversations.fulfilled, (state, action) => {
        state.loading = false
        state.conversations = action.payload
      })
      .addCase(fetchConversations.rejected, (state) => {
        state.loading = false
      })
      .addCase(fetchMessages.fulfilled, (state, action) => {
        const key = action.payload.type === 1 
          ? `private_${action.payload.targetId}` 
          : `group_${action.payload.targetId}`
        state.messages[key] = action.payload.messages
      })
  }
})

export const { setCurrentConversation, addMessage, setConversations, setMessages } = chatSlice.actions
export default chatSlice.reducer
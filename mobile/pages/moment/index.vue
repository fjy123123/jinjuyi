<template>
	<view class="moment-container">
		<!-- 发布栏 -->
		<view class="publish-bar" @click="goToPublish">
			<view class="publish-avatar">
				<image :src="userInfo.avatar || '/static/default-avatar.png'" mode="aspectFill" />
			</view>
			<view class="publish-input">说点什么...</view>
			<view class="publish-icon">
				<text class="iconfont icon-camera">📷</text>
			</view>
		</view>

		<!-- 朋友圈列表 -->
		<scroll-view scroll-y="true" @scrolltolower="loadMore" lower-threshold="100" :refreshing="refreshing" @refresherrefresh="onRefresh">
			<view v-for="moment in moments" :key="moment.id" class="moment-card">
				<!-- 头部 -->
				<view class="moment-header">
					<image :src="moment.user?.avatar || '/static/default-avatar.png'" mode="aspectFill" class="moment-avatar" />
					<view class="moment-user-info">
						<view class="moment-nickname">{{ moment.user?.nickname }}</view>
						<view class="moment-time">{{ formatTime(moment.created_at) }}</view>
					</view>
					<view class="moment-actions">
						<text @click="showActions(moment)" class="iconfont icon-more">...</text>
					</view>
				</view>

				<!-- 内容 -->
				<view v-if="moment.content" class="moment-content">{{ moment.content }}</view>

				<!-- 图片 -->
				<view v-if="moment.imagesArray && moment.imagesArray.length > 0" class="moment-images">
					<image 
						v-for="(img, idx) in moment.imagesArray" 
						:key="idx" 
						:src="img" 
						class="moment-img" 
						mode="aspectFill"
						@click="previewImages(idx, moment.imagesArray)"
					/>
				</view>

				<!-- 位置 -->
				<view v-if="moment.location" class="moment-location">
					<text class="iconfont icon-location">📍</text>
					<text>{{ moment.location }}</text>
				</view>

				<!-- 互动栏 -->
				<view class="moment-interact">
					<view class="interact-btn" @click="toggleLike(moment)">
						<text v-if="isLiked(moment)" class="liked">❤️</text>
						<text v-else>👍</text>
					</view>
					<view class="interact-btn" @click="focusComment(moment)">
						💬
					</view>
				</view>

				<!-- 点赞列表 -->
				<view v-if="moment.likes && moment.likes.length > 0" class="moment-likes">
					<text class="like-icon">❤️</text>
					<view class="like-users">
						<text 
							v-for="(like, idx) in moment.likes" 
							:key="like.id" 
							@click="goToUser(like.user_id)"
						>
							{{ like.user?.nickname }}
							<text v-if="idx < moment.likes.length - 1">, </text>
						</text>
					</view>
				</view>

				<!-- 评论列表 -->
				<view v-if="moment.comments && moment.comments.length > 0" class="moment-comments">
					<view v-for="comment in moment.comments" :key="comment.id" class="comment-item">
						<text class="comment-nick" @click="goToUser(comment.user_id)">{{ comment.user?.nickname }}</text>
						<text v-if="comment.reply_user" class="comment-reply">
							回复 <text @click="goToUser(comment.reply_to_user)">{{ comment.reply_user?.nickname }}</text>
						</text>
						<text>:</text>
						<text class="comment-content">{{ comment.content }}</text>
					</view>
				</view>
			</view>

			<!-- 加载更多 -->
			<view v-if="loading" class="loading-more">加载中...</view>
			<view v-if="noMore" class="no-more">没有更多了</view>
		</scroll-view>
	</view>
</template>

<script>
import { momentApi } from '@/api/moment'

export default {
	data() {
		return {
			userInfo: {},
			moments: [],
			page: 1,
			pageSize: 20,
			loading: false,
			noMore: false,
			refreshing: false
		}
	},
	onLoad() {
		this.loadMoments()
	},
	methods: {
		async loadMoments() {
			if (this.loading) return
			
			this.loading = true
			
			try {
				const res = await momentApi.getList(this.page, this.pageSize)
				
				if (this.page === 1) {
					this.moments = res.data.list
				} else {
					this.moments = [...this.moments, ...res.data.list]
				}
				
				if (res.data.list.length < this.pageSize) {
					this.noMore = true
				}
			} catch (err) {
				uni.showToast({
					title: '加载失败',
					icon: 'none'
				})
			} finally {
				this.loading = false
				this.refreshing = false
			}
		},
		loadMore() {
			if (!this.noMore && !this.loading) {
				this.page++
				this.loadMoments()
			}
		},
		onRefresh() {
			this.page = 1
			this.noMore = false
			this.loadMoments()
		},
		async toggleLike(moment) {
			if (this.isLiked(moment)) {
				await momentApi.unlike(moment.id)
				moment.likes = moment.likes.filter(l => l.user_id !== this.userInfo.id)
			} else {
				await momentApi.like(moment.id)
				// 刷新数据
				this.onRefresh()
			}
		},
		isLiked(moment) {
			return moment.likes?.some(l => l.user_id === this.userInfo.id)
		},
		focusComment(moment) {
			uni.showToast({
				title: '评论功能',
				icon: 'none'
			})
		},
		previewImages(index, images) {
			uni.previewImage({
				current: index,
				urls: images
			})
		},
		showActions(moment) {
			uni.showActionSheet({
				itemList: ['举报', '不看他', '取消'],
				success(res) {
					console.log(res.tapIndex)
				}
			})
		},
		goToPublish() {
			uni.navigateTo({
				url: '/pages/moment/publish'
			})
		},
		goToUser(userId) {
			uni.navigateTo({
				url: `/pages/user/index?id=${userId}`
			})
		},
		formatTime(time) {
			const now = Date.now()
			const createTime = new Date(time).getTime()
			const diff = now - createTime
			
			if (diff < 60000) {
				return '刚刚'
			} else if (diff < 3600000) {
				return Math.floor(diff / 60000) + '分钟前'
			} else if (diff < 86400000) {
				return Math.floor(diff / 3600000) + '小时前'
			} else if (diff < 2592000000) {
				return Math.floor(diff / 86400000) + '天前'
			} else {
				return time.substring(5, 16)
			}
		}
	}
}
</script>

<style scoped lang="scss">
.moment-container {
	padding-bottom: 30rpx;
	background-color: #f5f5f5;
	min-height: 100vh;
}

.publish-bar {
	display: flex;
	align-items: center;
	padding: 20rpx 30rpx;
	background-color: #fff;
	margin-bottom: 20rpx;
}

.publish-avatar image {
	width: 80rpx;
	height: 80rpx;
	border-radius: 8rpx;
}

.publish-input {
	flex: 1;
	margin-left: 20rpx;
	color: #999;
	font-size: 30rpx;
}

.publish-icon {
	font-size: 40rpx;
}

.moment-card {
	background-color: #fff;
	padding: 30rpx;
	margin-bottom: 20rpx;
}

.moment-header {
	display: flex;
	align-items: center;
}

.moment-avatar {
	width: 80rpx;
	height: 80rpx;
	border-radius: 8rpx;
}

.moment-user-info {
	flex: 1;
	margin-left: 20rpx;
}

.moment-nickname {
	font-size: 32rpx;
	font-weight: 500;
}

.moment-time {
	font-size: 24rpx;
	color: #999;
	margin-top: 6rpx;
}

.moment-content {
	margin-top: 20rpx;
	font-size: 30rpx;
	line-height: 1.6;
}

.moment-images {
	display: flex;
	flex-wrap: wrap;
	margin-top: 20rpx;
}

.moment-img {
	width: 200rpx;
	height: 200rpx;
	margin-right: 15rpx;
	margin-bottom: 15rpx;
	border-radius: 6rpx;
}

.moment-location {
	margin-top: 16rpx;
	font-size: 24rpx;
	color: #999;
}

.moment-interact {
	display: flex;
	justify-content: flex-end;
	margin-top: 20rpx;
	padding-top: 20rpx;
	border-top: 1rpx solid #eee;
}

.interact-btn {
	padding: 10rpx 30rpx;
	font-size: 30rpx;
}

.moment-likes {
	display: flex;
	align-items: center;
	background-color: #f7f7f7;
	padding: 15rpx 20rpx;
	margin-top: 15rpx;
	border-radius: 4rpx;
}

.like-users {
	flex: 1;
	margin-left: 10rpx;
	font-size: 28rpx;
	color: #333;
}

.moment-comments {
	background-color: #f7f7f7;
	padding: 15rpx 20rpx;
	margin-top: 10rpx;
	border-radius: 4rpx;
}

.comment-item {
	font-size: 28rpx;
	margin-bottom: 8rpx;
}

.comment-nick {
	color: #576b95;
}

.comment-reply {
	color: #576b95;
}

.loading-more, .no-more {
	text-align: center;
	padding: 30rpx;
	color: #999;
	font-size: 26rpx;
}
</style>

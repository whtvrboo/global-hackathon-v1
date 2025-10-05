<template>
    <div v-if="show"
        class="absolute z-50 bg-dark-900 border border-dark-700 rounded-xl shadow-2xl p-4 w-80 transform transition-all duration-200"
        :class="positionClass" @mouseenter="$emit('mouseenter')" @mouseleave="$emit('mouseleave')">
        <!-- User Info -->
        <div class="flex items-center gap-3 mb-4">
            <img v-if="user.picture" :src="user.picture" :alt="user.name"
                class="w-12 h-12 rounded-full border-2 border-dark-600" />
            <div v-else
                class="w-12 h-12 bg-dark-700 rounded-full flex items-center justify-center text-dark-300 border-2 border-dark-600 font-bold text-lg">
                {{ user.name.charAt(0).toUpperCase() }}
            </div>
            <div class="flex-1 min-w-0">
                <h4 class="font-bold text-white text-sm truncate">{{ user.name }}</h4>
                <p class="text-xs text-dark-400 truncate">@{{ user.username }}</p>
            </div>
        </div>

        <!-- Bio -->
        <div v-if="user.bio" class="mb-4">
            <p class="text-xs text-dark-200 leading-relaxed line-clamp-3">{{ user.bio }}</p>
        </div>

        <!-- Stats -->
        <div class="grid grid-cols-3 gap-3 mb-4">
            <div class="text-center">
                <div class="text-sm font-bold text-accent-blue">{{ user.stats?.followers_count || 0 }}</div>
                <div class="text-xs text-dark-400">Followers</div>
            </div>
            <div class="text-center">
                <div class="text-sm font-bold text-accent-green">{{ user.stats?.following_count || 0 }}</div>
                <div class="text-xs text-dark-400">Following</div>
            </div>
            <div class="text-center">
                <div class="text-sm font-bold text-accent-purple">{{ user.stats?.public_lists || 0 }}</div>
                <div class="text-xs text-dark-400">Lists</div>
            </div>
        </div>

        <!-- Actions -->
        <div class="flex gap-2">
            <button v-if="!isOwnProfile" @click="handleFollow" :disabled="followLoading"
                class="flex-1 px-3 py-2 text-xs font-medium rounded-lg transition-all duration-200" :class="isFollowing
                    ? 'bg-dark-700 text-white border border-dark-600 hover:bg-dark-600'
                    : 'bg-accent-red text-white hover:bg-accent-red/90'">
                {{ followLoading ? '...' : (isFollowing ? '✓ Following' : '+ Follow') }}
            </button>
            <button @click="viewProfile"
                class="flex-1 px-3 py-2 text-xs font-medium bg-dark-800 text-white rounded-lg border border-dark-600 hover:bg-dark-700 transition-all duration-200">
                View Profile
            </button>
        </div>
    </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import axios from 'axios'
import { useToastStore } from '../stores/toast'

const props = defineProps({
    show: Boolean,
    user: Object,
    position: {
        type: String,
        default: 'bottom-right' // bottom-right, bottom-left, top-right, top-left
    }
})

const emit = defineEmits(['close', 'mouseenter', 'mouseleave'])

const router = useRouter()
const authStore = useAuthStore()
const toastStore = useToastStore()

const isFollowing = ref(false)
const followLoading = ref(false)

const isOwnProfile = computed(() => {
    return authStore.user?.username === props.user?.username
})

const positionClass = computed(() => {
    const classes = {
        'bottom-right': 'top-full left-1/2 transform -translate-x-1/2 mt-2',
        'bottom-left': 'top-full right-1/2 transform translate-x-1/2 mt-2',
        'top-right': 'bottom-full left-1/2 transform -translate-x-1/2 mb-2',
        'top-left': 'bottom-full right-1/2 transform translate-x-1/2 mb-2'
    }
    return classes[props.position] || classes['bottom-right']
})


const handleFollow = async () => {
    if (followLoading.value || isOwnProfile.value) return

    followLoading.value = true
    try {
        if (isFollowing.value) {
            await axios.delete(`/api/users/${props.user.username}/follow`)
            isFollowing.value = false
            toastStore.info('Unfollowed')
        } else {
            await axios.post(`/api/users/${props.user.username}/follow`)
            isFollowing.value = true
            toastStore.success(`You're now following ${props.user.name}!`)
        }
    } catch (error) {
        console.error('Error toggling follow:', error)
        toastStore.error('Failed to update follow status')
    } finally {
        followLoading.value = false
    }
}

const viewProfile = () => {
    router.push(`/profile/${props.user.username}`)
    emit('close')
}
</script>

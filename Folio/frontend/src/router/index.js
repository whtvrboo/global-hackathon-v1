import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import AuthCallback from '../views/AuthCallback.vue'
import ProfileView from '../views/ProfileView.vue'
import FeedView from '../views/FeedView.vue'
import DiscoverView from '../views/DiscoverView.vue'
import { useAuthStore } from '../stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/profile',
      name: 'profile',
      component: ProfileView
    },
    {
      path: '/feed',
      name: 'feed',
      component: FeedView
    },
    {
      path: '/discover',
      name: 'discover',
      component: DiscoverView
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { public: true }
    },
    {
      path: '/auth/callback',
      name: 'auth-callback',
      component: AuthCallback,
      meta: { public: true }
    }
  ]
})

// Navigation guard for authentication
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  const isPublic = to.meta.public

  if (!isPublic && !authStore.isAuthenticated) {
    next('/login')
  } else if (to.path === '/login' && authStore.isAuthenticated) {
    next('/')
  } else {
    next()
  }
})

export default router


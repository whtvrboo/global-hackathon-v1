import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '../views/DashboardView.vue'
import LoginView from '../views/LoginView.vue'
import AuthCallback from '../views/AuthCallback.vue'
import ProfileView from '../views/ProfileView.vue'
import ListDetailView from '../views/ListDetailView.vue'
import BookDetailView from '../views/BookDetailView.vue'
import NotebookView from '../views/NotebookView.vue'
import ThreadView from '../views/ThreadView.vue'
import ListCreateView from '../views/ListCreateView.vue'
import { useAuthStore } from '../stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: DashboardView,
      meta: { public: true }
    },
    {
      path: '/profile/:username',
      name: 'profile',
      component: ProfileView,
      meta: { public: true }
    },
    {
      path: '/feed',
      redirect: '/'
    },
    {
      path: '/discover',
      redirect: '/'
    },
    {
      path: '/lists/:id',
      name: 'list-detail',
      component: ListDetailView,
      meta: { public: true }
    },
    {
      path: '/books/:id',
      name: 'book-detail',
      component: BookDetailView,
      meta: { public: true }
    },
    {
      path: '/notebook',
      name: 'notebook',
      component: NotebookView
    },
    {
      path: '/notebook/threads/:tag',
      name: 'thread',
      component: ThreadView
    },
    {
      path: '/lists/create',
      name: 'list-create',
      component: ListCreateView
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


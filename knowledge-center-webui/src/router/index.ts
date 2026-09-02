import { createRouter, createWebHistory } from 'vue-router'
import ImportView from '@/views/ImportView.vue'
import ListView from '@/views/ListView.vue'
import ScrapingView from '@/views/ScrapingView.vue'
import ImportDetailView from '@/views/ImportDetailView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/list' },
    { path: '/import', name: 'import', component: ImportView },
    { path: '/list', name: 'list', component: ListView },
    { path: '/list/:id', name: 'import-detail', component: ImportDetailView },
    { path: '/scraping', name: 'scraping', component: ScrapingView },
  ],
})

export default router

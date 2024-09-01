import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import PayView from "@/views/PayView.vue";
import UniView from "@/views/UniView.vue";
import LoanStartView from "@/views/LoanStartView.vue";
import LoanCheckView from "@/views/LoanCheckView.vue";
import IsurStartView from "@/views/IsurStartView.vue";
import IsurCheckView from "@/views/IsurCheckView.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'ecosys',
      component: HomeView
    },
      {
        path: '/home',
        name: 'home',
        component: HomeView
      },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue')
    },
    {
      path:'/ecosys/insuranceStart',
        name: 'insurance',
        component: IsurStartView
    },
      {
          path:'/ecosys/insuranceCheck',
            name: 'insuranceCheck',
            component: IsurCheckView
      },
      {
          path:'/ecosys/loanStart',
          name: 'loanStart',
          component: LoanStartView
      },
      {
          path:'/ecosys/loanCheck',
            name: 'loanCheck',
            component: LoanCheckView
      },
      {
          path:'/ecosys/pay',
          name: 'pay',
          component: PayView
      },
      {
          path:'/ecosys',
          name: 'universal',
          component: UniView
      }
  ]
})

export default router

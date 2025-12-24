import { createRouter, createWebHistory } from 'vue-router'
import BookList from '../views/BookList.vue'
import BookForm from '../views/BookForm.vue'
import LocationManage from '../views/LocationManage.vue'
import BorrowPage from '../views/BorrowPage.vue'
import BorrowQuery from '../views/BorrowQuery.vue'
import ReturnPage from '../views/ReturnPage.vue'

const routes = [
  {
    path: '/',
    redirect: '/books'
  },
  {
    path: '/books',
    name: 'BookList',
    component: BookList
  },
  {
    path: '/books/add',
    name: 'BookAdd',
    component: BookForm
  },
  {
    path: '/books/edit/:id',
    name: 'BookEdit',
    component: BookForm
  },
  {
    path: '/locations',
    name: 'LocationManage',
    component: LocationManage
  },
  {
    path: '/borrow',
    name: 'BorrowPage',
    component: BorrowPage
  },
  {
    path: '/borrow/query',
    name: 'BorrowQuery',
    component: BorrowQuery
  },
  {
    path: '/return',
    name: 'ReturnPage',
    component: ReturnPage
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router





import api from './index'

export const locationApi = {
  // 区域管理
  area: {
    create(data) {
      return api.post('/areas', data)
    },
    list() {
      return api.get('/areas')
    },
    update(id, data) {
      return api.put(`/areas/${id}`, data)
    },
    delete(id) {
      return api.delete(`/areas/${id}`)
    }
  },
  // 书架管理
  bookshelf: {
    create(data) {
      return api.post('/bookshelves', data)
    },
    list(areaId) {
      return api.get('/bookshelves', { params: { area_id: areaId } })
    },
    update(id, data) {
      return api.put(`/bookshelves/${id}`, data)
    },
    delete(id) {
      return api.delete(`/bookshelves/${id}`)
    }
  },
  // 层数管理
  shelfLayer: {
    create(data) {
      return api.post('/shelf-layers', data)
    },
    list(bookshelfId) {
      return api.get('/shelf-layers', { params: { bookshelf_id: bookshelfId } })
    },
    update(id, data) {
      return api.put(`/shelf-layers/${id}`, data)
    },
    delete(id) {
      return api.delete(`/shelf-layers/${id}`)
    }
  },
  // 获取位置树
  getTree() {
    return api.get('/locations/tree')
  }
}




import api from './index'

export const bookApi = {
  // 创建图书
  create(data) {
    return api.post('/books', data)
  },
  // 查询图书列表
  list(params) {
    return api.get('/books', { params })
  },
  // 根据一维码查询
  getByBarcode(barcode) {
    return api.get(`/books/barcode/${barcode}`)
  },
  // 更新图书
  update(id, data) {
    return api.put(`/books/${id}`, data)
  },
  // 删除图书
  delete(id) {
    return api.delete(`/books/${id}`)
  }
}




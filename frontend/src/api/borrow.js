import api from './index'

export const borrowApi = {
  // 创建借阅记录（兼容旧接口）
  create(data) {
    return api.post('/borrow', data)
  },
  // 扫码借阅（单本）
  scan(data) {
    return api.post('/borrow/scan', data)
  },
  // 查询借阅记录
  list(params) {
    return api.get('/borrow/records', { params })
  },
  // 根据电话查询归还人信息
  getBorrowerByPhone(data) {
    return api.post('/borrow/get-borrower', data)
  },
  // 归还（兼容旧接口）
  returnBook(data) {
    return api.post('/borrow/return', data)
  },
  // 新的借阅API（使用Redis）
  setBorrowUser(data) {
    return api.post('/borrow/user', data)
  },
  getBorrowUser() {
    return api.get('/borrow/user')
  },
  addBorrowBook(data) {
    return api.post('/borrow/book', data)
  },
  removeBorrowBook(data) {
    return api.delete('/borrow/book', { data: data })
  },
  completeBorrow(data) {
    return api.post('/borrow/complete', data)
  },
  // 新的归还API（使用Redis）
  setReturnUser(data) {
    return api.post('/return/user', data)
  },
  getReturnUser() {
    return api.get('/return/user')
  },
  addReturnBook(data) {
    return api.post('/return/book', data)
  },
  removeReturnBook(data) {
    return api.delete('/return/book', { data: data })
  },
  completeReturn(data) {
    return api.post('/return/complete', data)
  }
}



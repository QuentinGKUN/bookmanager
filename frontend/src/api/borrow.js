import api from './index'

export const borrowApi = {
  // 创建借阅记录
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
  // 归还
  returnBook(data) {
    return api.post('/borrow/return', data)
  }
}



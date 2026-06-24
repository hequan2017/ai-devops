import service from '@/utils/request'

// 发版 CRUD
export const createRelease = (data) => service({ url: '/release/createRelease', method: 'post', data })
export const deleteRelease = (data) => service({ url: '/release/deleteRelease', method: 'delete', data })
export const deleteReleaseByIds = (data) => service({ url: '/release/deleteReleaseByIds', method: 'delete', data })
export const updateRelease = (data) => service({ url: '/release/updateRelease', method: 'put', data })
export const findRelease = (params) => service({ url: '/release/findRelease', method: 'get', params })
export const getReleaseList = (params) => service({ url: '/release/getReleaseList', method: 'get', params })

// 发版工作流
export const submitRelease = (data) => service({ url: '/release/submitRelease', method: 'post', data })
export const approveRelease = (data) => service({ url: '/release/approveRelease', method: 'post', data })
export const rejectRelease = (data) => service({ url: '/release/rejectRelease', method: 'post', data })
export const executeRelease = (data) => service({ url: '/release/executeRelease', method: 'post', data })

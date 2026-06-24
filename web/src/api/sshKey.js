import service from '@/utils/request'

export const createSshKey = (data) => service({ url: '/sshKey/createSshKey', method: 'post', data })
export const deleteSshKey = (data) => service({ url: '/sshKey/deleteSshKey', method: 'delete', data })
export const deleteSshKeyByIds = (data) => service({ url: '/sshKey/deleteSshKeyByIds', method: 'delete', data })
export const updateSshKey = (data) => service({ url: '/sshKey/updateSshKey', method: 'put', data })
export const findSshKey = (params) => service({ url: '/sshKey/findSshKey', method: 'get', params })
export const getSshKeyList = (params) => service({ url: '/sshKey/getSshKeyList', method: 'get', params })
export const getAllSshKey = () => service({ url: '/sshKey/getAllSshKey', method: 'get' })

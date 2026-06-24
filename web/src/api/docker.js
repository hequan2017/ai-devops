import service from '@/utils/request'

// Docker 接入点 CRUD
export const createDockerHost = (data) => service({ url: '/docker/createDockerHost', method: 'post', data })
export const deleteDockerHost = (data) => service({ url: '/docker/deleteDockerHost', method: 'delete', data })
export const deleteDockerHostByIds = (data) => service({ url: '/docker/deleteDockerHostByIds', method: 'delete', data })
export const updateDockerHost = (data) => service({ url: '/docker/updateDockerHost', method: 'put', data })
export const findDockerHost = (params) => service({ url: '/docker/findDockerHost', method: 'get', params })
export const getDockerHostList = (params) => service({ url: '/docker/getDockerHostList', method: 'get', params })

// Docker 操作
export const testDockerHost = (params) => service({ url: '/docker/testDockerHost', method: 'get', params })
export const getContainers = (params) => service({ url: '/docker/getContainers', method: 'get', params })
export const containerAction = (data) => service({ url: '/docker/containerAction', method: 'post', data })
export const getContainerLogs = (params) => service({ url: '/docker/getContainerLogs', method: 'get', params })
export const getImages = (params) => service({ url: '/docker/getImages', method: 'get', params })
export const getNetworks = (params) => service({ url: '/docker/getNetworks', method: 'get', params })
export const getVolumes = (params) => service({ url: '/docker/getVolumes', method: 'get', params })

import service from '@/utils/request'

// K8s 集群 CRUD
export const createK8sCluster = (data) => service({ url: '/k8s/createK8sCluster', method: 'post', data })
export const deleteK8sCluster = (data) => service({ url: '/k8s/deleteK8sCluster', method: 'delete', data })
export const deleteK8sClusterByIds = (data) => service({ url: '/k8s/deleteK8sClusterByIds', method: 'delete', data })
export const updateK8sCluster = (data) => service({ url: '/k8s/updateK8sCluster', method: 'put', data })
export const findK8sCluster = (params) => service({ url: '/k8s/findK8sCluster', method: 'get', params })
export const getK8sClusterList = (params) => service({ url: '/k8s/getK8sClusterList', method: 'get', params })

// K8s 操作
export const testK8sCluster = (params) => service({ url: '/k8s/testK8sCluster', method: 'get', params })
export const getNamespaces = (params) => service({ url: '/k8s/getNamespaces', method: 'get', params })
export const getPods = (params) => service({ url: '/k8s/getPods', method: 'get', params })
export const getNodes = (params) => service({ url: '/k8s/getNodes', method: 'get', params })
export const getDeployments = (params) => service({ url: '/k8s/getDeployments', method: 'get', params })
export const getServices = (params) => service({ url: '/k8s/getServices', method: 'get', params })
export const getPodLogs = (params) => service({ url: '/k8s/getPodLogs', method: 'get', params })
export const deletePod = (data) => service({ url: '/k8s/deletePod', method: 'post', data })

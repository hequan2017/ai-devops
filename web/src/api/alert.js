import service from '@/utils/request'

// 告警规则
export const createAlertRule = (data) => service({ url: '/alert/createAlertRule', method: 'post', data })
export const deleteAlertRule = (data) => service({ url: '/alert/deleteAlertRule', method: 'delete', data })
export const deleteAlertRuleByIds = (data) => service({ url: '/alert/deleteAlertRuleByIds', method: 'delete', data })
export const updateAlertRule = (data) => service({ url: '/alert/updateAlertRule', method: 'put', data })
export const findAlertRule = (params) => service({ url: '/alert/findAlertRule', method: 'get', params })
export const getAlertRuleList = (params) => service({ url: '/alert/getAlertRuleList', method: 'get', params })

// 告警记录
export const getAlertRecordList = (params) => service({ url: '/alert/getAlertRecordList', method: 'get', params })
export const resolveAlertRecord = (data) => service({ url: '/alert/resolveAlertRecord', method: 'post', data })

import service from '@/utils/request'

// @Tags Server
// @Summary 新增服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body biz.Server true "新增服务器"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /server/createServer [post]
export const createServer = (data) => {
  return service({
    url: '/server/createServer',
    method: 'post',
    data
  })
}

// @Tags Server
// @Summary 删除服务器
// @Router /server/deleteServer [delete]
export const deleteServer = (data) => {
  return service({
    url: '/server/deleteServer',
    method: 'delete',
    data
  })
}

// @Tags Server
// @Summary 批量删除服务器
// @Router /server/deleteServerByIds [delete]
export const deleteServerByIds = (data) => {
  return service({
    url: '/server/deleteServerByIds',
    method: 'delete',
    data
  })
}

// @Tags Server
// @Summary 更新服务器
// @Router /server/updateServer [put]
export const updateServer = (data) => {
  return service({
    url: '/server/updateServer',
    method: 'put',
    data
  })
}

// @Tags Server
// @Summary 根据ID查询服务器
// @Router /server/findServer [get]
export const findServer = (params) => {
  return service({
    url: '/server/findServer',
    method: 'get',
    params
  })
}

// @Tags Server
// @Summary 分页获取服务器列表
// @Router /server/getServerList [get]
export const getServerList = (params) => {
  return service({
    url: '/server/getServerList',
    method: 'get',
    params
  })
}

// @Tags Server
// @Summary IPMI 电源控制
// @Router /server/ipmiPower [post]
export const ipmiPower = (data) => {
  return service({
    url: '/server/ipmiPower',
    method: 'post',
    data
  })
}

// @Tags Server
// @Summary 批量执行命令
// @Router /server/execCmd [post]
export const execCmd = (data) => {
  return service({
    url: '/server/execCmd',
    method: 'post',
    data
  })
}

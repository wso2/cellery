/* eslint-disable */
import axios from 'axios'
import qs from 'qs'
let domain = ''
export const getDomain = () => {
  return domain
}
export const setDomain = ($domain) => {
  domain = $domain
}
export const request = (method, url, body, queryParameters, form, config) => {
  method = method.toLowerCase()
  let keys = Object.keys(queryParameters)
  let queryUrl = url
  if (keys.length > 0) {
    queryUrl = url + '?' + qs.stringify(queryParameters)
  }
  // let queryUrl = url+(keys.length > 0 ? '?' + (keys.map(key => key + '=' + encodeURIComponent(queryParameters[key])).join('&')) : '')
  if (body) {
    return axios[method](queryUrl, body, config)
  } else if (method === 'get' || method === 'delete' || method === 'head' || method === 'option') {
    return axios[method](queryUrl, config)
  } else {
    return axios[method](queryUrl, qs.stringify(form), config)
  }
}
/*==========================================================
 *                    Cellery sample server Petstore server.
 ==========================================================*/
/**
 * View all pet accessories in the store
 * request: getCatalog
 * url: getCatalogURL
 * method: getCatalog_TYPE
 * raw_url: getCatalog_RAW_URL
 */
export const getCatalog = function(parameters = {}) {
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  const config = parameters.$config || {
    headers: {}
  }
  let path = '/catalog'
  let body
  let queryParameters = {}
  let form = {}
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    });
  }
  return request('get', domain + path, body, queryParameters, form, config)
}
export const getCatalog_RAW_URL = function() {
  return '/catalog'
}
export const getCatalog_TYPE = function() {
  return 'get'
}
export const getCatalogURL = function(parameters = {}) {
  let queryParameters = {}
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  let path = '/catalog'
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    })
  }
  let keys = Object.keys(queryParameters)
  return domain + path + (keys.length > 0 ? '?' + (keys.map(key => key + '=' + encodeURIComponent(queryParameters[key])).join('&')) : '')
}
/**
 * View all pet accessories in the store
 * request: getProduct
 * url: getProductURL
 * method: getProduct_TYPE
 * raw_url: getProduct_RAW_URL
 * @param productId - ID of product to return
 */
export const getProduct = function(parameters = {}) {
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  const config = parameters.$config || {
    headers: {}
  }
  let path = '/catalog/{product_id}'
  let body
  let queryParameters = {}
  let form = {}
  path = path.replace('{product_id}', `${parameters['productId']}`)
  if (parameters['productId'] === undefined) {
    return Promise.reject(new Error('Missing required  parameter: productId'))
  }
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    });
  }
  return request('get', domain + path, body, queryParameters, form, config)
}
export const getProduct_RAW_URL = function() {
  return '/catalog/{product_id}'
}
export const getProduct_TYPE = function() {
  return 'get'
}
export const getProductURL = function(parameters = {}) {
  let queryParameters = {}
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  let path = '/catalog/{product_id}'
  path = path.replace('{product_id}', `${parameters['productId']}`)
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    })
  }
  let keys = Object.keys(queryParameters)
  return domain + path + (keys.length > 0 ? '?' + (keys.map(key => key + '=' + encodeURIComponent(queryParameters[key])).join('&')) : '')
}
/**
 * View all customers
 * request: getCustomers
 * url: getCustomersURL
 * method: getCustomers_TYPE
 * raw_url: getCustomers_RAW_URL
 */
export const getCustomers = function(parameters = {}) {
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  const config = parameters.$config || {
    headers: {}
  }
  let path = '/customers'
  let body
  let queryParameters = {}
  let form = {}
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    });
  }
  return request('get', domain + path, body, queryParameters, form, config)
}
export const getCustomers_RAW_URL = function() {
  return '/customers'
}
export const getCustomers_TYPE = function() {
  return 'get'
}
export const getCustomersURL = function(parameters = {}) {
  let queryParameters = {}
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  let path = '/customers'
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    })
  }
  let keys = Object.keys(queryParameters)
  return domain + path + (keys.length > 0 ? '?' + (keys.map(key => key + '=' + encodeURIComponent(queryParameters[key])).join('&')) : '')
}
/**
 * View all customers
 * request: getCustomer
 * url: getCustomerURL
 * method: getCustomer_TYPE
 * raw_url: getCustomer_RAW_URL
 * @param customerId - ID of customer to return
 */
export const getCustomer = function(parameters = {}) {
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  const config = parameters.$config || {
    headers: {}
  }
  let path = '/customers/{customer_id}'
  let body
  let queryParameters = {}
  let form = {}
  path = path.replace('{customer_id}', `${parameters['customerId']}`)
  if (parameters['customerId'] === undefined) {
    return Promise.reject(new Error('Missing required  parameter: customerId'))
  }
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    });
  }
  return request('get', domain + path, body, queryParameters, form, config)
}
export const getCustomer_RAW_URL = function() {
  return '/customers/{customer_id}'
}
export const getCustomer_TYPE = function() {
  return 'get'
}
export const getCustomerURL = function(parameters = {}) {
  let queryParameters = {}
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  let path = '/customers/{customer_id}'
  path = path.replace('{customer_id}', `${parameters['customerId']}`)
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    })
  }
  let keys = Object.keys(queryParameters)
  return domain + path + (keys.length > 0 ? '?' + (keys.map(key => key + '=' + encodeURIComponent(queryParameters[key])).join('&')) : '')
}
/**
 * View all orders
 * request: getOrders
 * url: getOrdersURL
 * method: getOrders_TYPE
 * raw_url: getOrders_RAW_URL
 */
export const getOrders = function(parameters = {}) {
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  const config = parameters.$config || {
    headers: {}
  }
  let path = '/orders'
  let body
  let queryParameters = {}
  let form = {}
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    });
  }
  return request('get', domain + path, body, queryParameters, form, config)
}
export const getOrders_RAW_URL = function() {
  return '/orders'
}
export const getOrders_TYPE = function() {
  return 'get'
}
export const getOrdersURL = function(parameters = {}) {
  let queryParameters = {}
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  let path = '/orders'
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    })
  }
  let keys = Object.keys(queryParameters)
  return domain + path + (keys.length > 0 ? '?' + (keys.map(key => key + '=' + encodeURIComponent(queryParameters[key])).join('&')) : '')
}
/**
 * View all orders
 * request: getOrder
 * url: getOrderURL
 * method: getOrder_TYPE
 * raw_url: getOrder_RAW_URL
 * @param orderId - ID of customer to return
 */
export const getOrder = function(parameters = {}) {
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  const config = parameters.$config || {
    headers: {}
  }
  let path = '/orders/{order_id}'
  let body
  let queryParameters = {}
  let form = {}
  path = path.replace('{order_id}', `${parameters['orderId']}`)
  if (parameters['orderId'] === undefined) {
    return Promise.reject(new Error('Missing required  parameter: orderId'))
  }
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    });
  }
  return request('get', domain + path, body, queryParameters, form, config)
}
export const getOrder_RAW_URL = function() {
  return '/orders/{order_id}'
}
export const getOrder_TYPE = function() {
  return 'get'
}
export const getOrderURL = function(parameters = {}) {
  let queryParameters = {}
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  let path = '/orders/{order_id}'
  path = path.replace('{order_id}', `${parameters['orderId']}`)
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    })
  }
  let keys = Object.keys(queryParameters)
  return domain + path + (keys.length > 0 ? '?' + (keys.map(key => key + '=' + encodeURIComponent(queryParameters[key])).join('&')) : '')
}
/**
 * View all orders of a specific customer
 * request: getOrderByCustomer
 * url: getOrderByCustomerURL
 * method: getOrderByCustomer_TYPE
 * raw_url: getOrderByCustomer_RAW_URL
 * @param customerId - ID of customer to search orders for
 */
export const getOrderByCustomer = function(parameters = {}) {
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  const config = parameters.$config || {
    headers: {}
  }
  let path = '/orders/customer/{customer_id}'
  let body
  let queryParameters = {}
  let form = {}
  path = path.replace('{customer_id}', `${parameters['customerId']}`)
  if (parameters['customerId'] === undefined) {
    return Promise.reject(new Error('Missing required  parameter: customerId'))
  }
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    });
  }
  return request('get', domain + path, body, queryParameters, form, config)
}
export const getOrderByCustomer_RAW_URL = function() {
  return '/orders/customer/{customer_id}'
}
export const getOrderByCustomer_TYPE = function() {
  return 'get'
}
export const getOrderByCustomerURL = function(parameters = {}) {
  let queryParameters = {}
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  let path = '/orders/customer/{customer_id}'
  path = path.replace('{customer_id}', `${parameters['customerId']}`)
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    })
  }
  let keys = Object.keys(queryParameters)
  return domain + path + (keys.length > 0 ? '?' + (keys.map(key => key + '=' + encodeURIComponent(queryParameters[key])).join('&')) : '')
}
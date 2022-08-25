// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "https://www.wzxaugenstern.online/#/",
        "contact": {
            "name": "CYCLEWW",
            "url": "https://www.wzxaugenstern.online/#/",
            "email": "1539989223@qq.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/message/election": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "由领导者调用的新一轮选举接口。",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "服务层"
                ],
                "summary": "选举",
                "parameters": [
                    {
                        "type": "string",
                        "name": "IP",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "colony",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "namespace",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "port",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "time",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "返回true或false",
                        "schema": {
                            "$ref": "#/definitions/entity.ResultModel"
                        }
                    }
                }
            }
        },
        "/api/v1/message/getLeader": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "由其他服务调用的获取当前领导者接口。",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "服务层"
                ],
                "summary": "获取领导者",
                "parameters": [
                    {
                        "type": "string",
                        "name": "IP",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "colony",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "namespace",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "port",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "time",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "返回领导者服务信息",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/entity.ResultModel"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/entity.ServerModel"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/message/getServers": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "由当前领导者调用的获取领导者所领导的服务列表。",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "服务层"
                ],
                "summary": "获取当前被领导者服务列表",
                "parameters": [
                    {
                        "type": "string",
                        "name": "IP",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "colony",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "namespace",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "port",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "time",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "返回被领导者的切片数组",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/entity.ResultModel"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/entity.ServerModel"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/message/message/beat": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "服务心跳重置倒计时",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "服务层"
                ],
                "summary": "服务心跳",
                "parameters": [
                    {
                        "type": "string",
                        "name": "IP",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "colony",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "namespace",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "port",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "time",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "返回true或false",
                        "schema": {
                            "$ref": "#/definitions/entity.ResultModel"
                        }
                    }
                }
            }
        },
        "/api/v1/message/register": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "服务注册进中心",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "服务层"
                ],
                "summary": "服务注册",
                "parameters": [
                    {
                        "type": "string",
                        "name": "IP",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "colony",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "namespace",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "port",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "time",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "返回true或false",
                        "schema": {
                            "$ref": "#/definitions/entity.ResultModel"
                        }
                    }
                }
            }
        },
        "/api/v1/operator/cancelDeleteInstance": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "由管理者调用删除黑名单中的实例信息。",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "管理层"
                ],
                "summary": "删除黑名单中的实例信息",
                "parameters": [
                    {
                        "type": "string",
                        "name": "IP",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "colony",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "namespace",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "port",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "time",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "返回true或false",
                        "schema": {
                            "$ref": "#/definitions/entity.ResultModel"
                        }
                    }
                }
            }
        },
        "/api/v1/operator/deleteColony": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "由管理者调用删除地区集群实例并拉入黑名单。",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "管理层"
                ],
                "summary": "删除地区集群实例并拉入黑名单",
                "parameters": [
                    {
                        "type": "string",
                        "name": "IP",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "colony",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "namespace",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "port",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "time",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "返回true或false",
                        "schema": {
                            "$ref": "#/definitions/entity.ResultModel"
                        }
                    }
                }
            }
        },
        "/api/v1/operator/election": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "由管理者调用删除服务实例并拉入黑名单。",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "管理层"
                ],
                "summary": "删除服务实例并拉入黑名单",
                "parameters": [
                    {
                        "type": "string",
                        "name": "IP",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "colony",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "namespace",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "port",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "time",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "返回true或false",
                        "schema": {
                            "$ref": "#/definitions/entity.ResultModel"
                        }
                    }
                }
            }
        },
        "/api/v1/operator/getDeleteInstance": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "由管理者调用的获取当前全部黑名单服务实例",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "管理层"
                ],
                "summary": "获取全部黑名单服务实例",
                "responses": {
                    "200": {
                        "description": "返回黑名单中服务实例切片数组",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/entity.ResultModel"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/entity.ServerModel"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/operator/getInstances": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "由管理者调用的获取当前所有服务信息。",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "管理层"
                ],
                "summary": "获取全部服务实例",
                "responses": {
                    "200": {
                        "description": "返回服务实例切片数组",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/entity.ResultModel"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/entity.ServerModel"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "由获取命名空间与集群条件下的服务实例服务实例。",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "管理层"
                ],
                "summary": "获取指定条件下的服务实例服务实例",
                "responses": {
                    "200": {
                        "description": "返回服务实例切片数组",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/entity.ResultModel"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/entity.ServerModel"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/operator/getStatus": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "由管理员调用获取当前中心线程数",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "管理层"
                ],
                "summary": "获取服务状态",
                "responses": {
                    "200": {
                        "description": "返回电脑状态",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/entity.ResultModel"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/entity.ComputerInfoModel"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.ComputerInfoModel": {
            "type": "object",
            "properties": {
                "cpu_info": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.CpuInfoModel"
                    }
                },
                "disk_info": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.DiskInfoModel"
                    }
                },
                "host_info": {
                    "$ref": "#/definitions/entity.HostInfoModel"
                },
                "mem_info": {
                    "$ref": "#/definitions/entity.MemInfoModel"
                },
                "net_info": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.NetInfoModel"
                    }
                },
                "pool_activate_num": {
                    "type": "integer"
                },
                "pool_job_num": {
                    "type": "integer"
                }
            }
        },
        "entity.CpuInfoModel": {
            "type": "object",
            "properties": {
                "cpu_core_num": {
                    "type": "string"
                },
                "cpu_frequency": {
                    "type": "string"
                },
                "cpu_name": {
                    "type": "string"
                },
                "cpu_physical_id": {
                    "type": "string"
                },
                "cpu_usage": {
                    "type": "string"
                },
                "cpu_vendor_id": {
                    "type": "string"
                }
            }
        },
        "entity.DiskInfoModel": {
            "type": "object",
            "properties": {
                "disk_avi": {
                    "type": "string"
                },
                "disk_name": {
                    "type": "string"
                },
                "disk_size": {
                    "type": "string"
                },
                "disk_usage": {
                    "type": "string"
                },
                "disk_used": {
                    "type": "string"
                },
                "fs_type": {
                    "type": "string"
                },
                "opts": {
                    "type": "string"
                }
            }
        },
        "entity.HostInfoModel": {
            "type": "object",
            "properties": {
                "host_id": {
                    "type": "string"
                },
                "host_kernel_arch": {
                    "type": "string"
                },
                "host_kernel_version": {
                    "type": "string"
                },
                "host_name": {
                    "type": "string"
                },
                "host_os": {
                    "type": "string"
                },
                "host_os_version": {
                    "type": "string"
                }
            }
        },
        "entity.MemInfoModel": {
            "type": "object",
            "properties": {
                "mem_avi": {
                    "type": "string"
                },
                "mem_total": {
                    "type": "string"
                },
                "mem_usage": {
                    "type": "string"
                },
                "mem_used": {
                    "type": "string"
                }
            }
        },
        "entity.NetInfoModel": {
            "type": "object",
            "properties": {
                "bytes_received": {
                    "type": "string"
                },
                "bytes_sent": {
                    "type": "string"
                },
                "net_name": {
                    "type": "string"
                },
                "packets_received": {
                    "type": "string"
                },
                "packets_sent": {
                    "type": "string"
                }
            }
        },
        "entity.ResultModel": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "message": {
                    "type": "string"
                }
            }
        },
        "entity.ServerModel": {
            "type": "object",
            "properties": {
                "IP": {
                    "type": "string"
                },
                "colony": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "namespace": {
                    "type": "string"
                },
                "port": {
                    "type": "string"
                },
                "time": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Themis API",
	Description:      "分布式记账系统调度中心",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

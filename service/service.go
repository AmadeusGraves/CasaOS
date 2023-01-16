/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-07-12 09:48:56
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-09-02 22:10:05
 * @FilePath: /CasaOS/service/service.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package service

import (
	"github.com/IceWhaleTech/CasaOS-Common/external"
	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/websocket"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var Cache *cache.Cache

var MyService Repository
var SocketServer *socketio.Server
var (
	WebSocketConns []*websocket.Conn
	SocketRun      bool
)

type Repository interface {
	// User() UserService
	Casa() CasaService
	Notify() NotifyServer
	Rely() RelyService
	System() SystemService
	Shares() SharesService
	Connections() ConnectionsService
	Gateway() external.ManagementService
	Storage() StorageService
	Storages() StoragesService
	StoragePath() StoragePathService
	FsListService() FsListService
	FsService() FsService
}

func NewService(db *gorm.DB, RuntimePath string, socket *socketio.Server) Repository {
	if socket == nil {
		logger.Error("socket is nil", zap.Any("error", "socket is nil"))
	}
	SocketServer = socket
	gatewayManagement, err := external.NewManagementService(RuntimePath)
	if err != nil && len(RuntimePath) > 0 {
		panic(err)
	}

	return &store{
		gateway:      gatewayManagement,
		casa:         NewCasaService(),
		notify:       NewNotifyService(db),
		rely:         NewRelyService(db),
		system:       NewSystemService(),
		shares:       NewSharesService(db),
		connections:  NewConnectionsService(db),
		storage:      NewStorageService(db),
		storages:     NewStoragesService(),
		storage_path: NewStoragePathService(),
		fs_list:      NewFsListService(),
		fs:           NewFsService(),
	}
}

type store struct {
	db           *gorm.DB
	casa         CasaService
	notify       NotifyServer
	rely         RelyService
	system       SystemService
	shares       SharesService
	connections  ConnectionsService
	gateway      external.ManagementService
	storage      StorageService
	storages     StoragesService
	storage_path StoragePathService
	fs_list      FsListService
	fs           FsService
}

func (c *store) FsService() FsService {
	return c.fs
}

func (c *store) FsListService() FsListService {
	return c.fs_list
}

func (c *store) StoragePath() StoragePathService {
	return c.storage_path
}

func (c *store) Storages() StoragesService {
	return c.storages
}

func (c *store) Storage() StorageService {
	return c.storage
}

func (c *store) Gateway() external.ManagementService {
	return c.gateway
}

func (s *store) Connections() ConnectionsService {
	return s.connections
}

func (s *store) Shares() SharesService {
	return s.shares
}

func (c *store) Rely() RelyService {
	return c.rely
}

func (c *store) System() SystemService {
	return c.system
}

func (c *store) Notify() NotifyServer {
	return c.notify
}

func (c *store) Casa() CasaService {
	return c.casa
}

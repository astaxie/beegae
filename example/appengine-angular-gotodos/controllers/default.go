// EXAMPLE FROM: https://github.com/GoogleCloudPlatform/appengine-angular-gotodos
//
// Copyright 2013 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// gotodos is an App Engine JSON backend for managing a todo list.
//
// It supports the following commands:
//
// - Create a new todo
// POST /todos
// > {"text": "do this"}
// < {"id": 1, "text": "do this", "created": 1356724843.0, "done": false}
//
// - Update an existing todo
// POST /todos
// > {"id": 1, "text": "do this", "created": 1356724843.0, "done": true}
// < {"id": 1, "text": "do this", "created": 1356724843.0, "done": true}
//
// - List existing todos:
// GET /todos
// >
// < [{"id": 1, "text": "do this", "created": 1356724843.0, "done": true},
//    {"id": 2, "text": "do that", "created": 1356724849.0, "done": false}]
//
// - Delete 'done' todos:
// DELETE /todos
// >
// <

package controllers

import (
	"encoding/json"
	"io"

	"models"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	"github.com/astaxie/beegae"
)

type MainController struct {
	beegae.Controller
}

func (this *MainController) Get() {
	todos := []models.Todo{}
	ks, err := datastore.NewQuery("Todo").Ancestor(models.DefaultTodoList(this.AppEngineCtx)).Order("Created").GetAll(this.AppEngineCtx, &todos)
	if err != nil {
		this.Data["json"] = err
		return
	}
	for i := 0; i < len(todos); i++ {
		todos[i].Id = ks[i].IntID()
	}
	this.Data["json"] = todos
}

func (this *MainController) Post() {
	todo, err := decodeTodo(this.Ctx.Input.Request.Body)
	if err != nil {
		this.Data["json"] = err
		return
	}
	t, err := todo.Save(this.AppEngineCtx)
	if err != nil {
		this.Data["json"] = err
	} else {
		this.Data["json"] = &t
	}
}

func (this *MainController) Delete() {
	err := datastore.RunInTransaction(this.AppEngineCtx, func(c context.Context) error {
		ks, err := datastore.NewQuery("Todo").KeysOnly().Ancestor(models.DefaultTodoList(c)).Filter("Done=", true).GetAll(c, nil)
		if err != nil {
			return err
		}
		return datastore.DeleteMulti(c, ks)
	}, nil)

	if err == nil {
		this.Data["json"] = nil
	} else {
		this.Data["json"] = err
	}
}

func (this *MainController) Render() error {
	if _, ok := this.Data["json"].(error); ok {
		log.Errorf(this.AppEngineCtx, "todo error: %v", this.Data["json"])
	}
	this.ServeJson()
	return nil
}

func decodeTodo(r io.ReadCloser) (*models.Todo, error) {
	defer r.Close()
	var todo models.Todo
	err := json.NewDecoder(r).Decode(&todo)
	return &todo, err
}

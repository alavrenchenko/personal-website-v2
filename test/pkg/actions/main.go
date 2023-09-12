// Copyright 2023 Alexey Lavrenchenko. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"

	"github.com/google/uuid"

	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/actions/logging"
	"personal-website-v2/pkg/logging/logger"
)

const (
	appSessionId uint64 = 1
)

func main() {
	f, err := logger.NewLoggerFactory(loggingSessionId, createLoggerConfig(), true)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := f.Dispose(); err != nil {
			fmt.Println(err)
		}
	}()

	l, err := logging.NewLogger(appSessionId, createActionLoggerConfig())

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := l.Dispose(); err != nil {
			fmt.Println(err)
		}
	}()

	tranManager, err := actions.NewTransactionManager(appSessionId, l, f)

	if err != nil {
		panic(err)
	}

	testTransactionManager(tranManager)

	actionManager, err := actions.NewActionManager(appSessionId, l, l, f)

	if err != nil {
		panic(err)
	}

	fmt.Println()
	testActionManager(tranManager, actionManager)

	fmt.Println()
	testOperationManager(tranManager, actionManager)

}

func testTransactionManager(m *actions.TransactionManager) {
	fmt.Println("***** testTransactionManager *****")

	for i := 0; i < 10; i++ {
		t, err := m.CreateAndStart()

		if err != nil {
			panic(err)
		}

		fmt.Printf("Transaction:\nId: %s\nCreatedAt: %s\nStartTime: %s\n\n", t.Id(), t.CreatedAt(), t.StartTime())
	}

	fmt.Printf("TransactionManager:\nCounter: %d\nNumCreated: %d\n", m.Counter(), m.NumCreated())
}

func testActionManager(tranManager *actions.TransactionManager, actionManager *actions.ActionManager) {
	fmt.Println("***** testActionManager *****")
	t, err := tranManager.CreateAndStart()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Transaction:\nId: %s\nCreatedAt: %s\nStartTime: %s\n\n", t.Id(), t.CreatedAt(), t.StartTime())

	for i := 1; i <= 10; i++ {
		a, err := actionManager.CreateAndStart(t, actions.ActionTypeApplication_Start, actions.ActionCategoryCommon, actions.ActionGroupApplication, uuid.NullUUID{}, false)

		if err != nil {
			panic(err)
		}

		fmt.Printf("Action:\nId: %s\nTranId: %s\nType: %v\nCategory: %v\nGroup: %v\nParentActionId: %v\nIsBackground: %v\nCreatedAt: %s\nStatus: %v\nStartTime: %s\nEndTime: %v\nElapsedTime: %v\nIsCompleted: %v\n\n",
			a.Id(), a.Transaction().Id(), a.Type(), a.Category(), a.Group(), a.ParentActionId(), a.IsBackground(), a.CreatedAt(), a.Status(), a.StartTime(), a.EndTime(), a.ElapsedTime(), a.IsCompleted())

		if err = actionManager.Complete(a, i%2 == 1); err != nil {
			panic(err)
		}

		fmt.Printf("Action:\nId: %s\nTranId: %s\nType: %v\nCategory: %v\nGroup: %v\nParentActionId: %v\nIsBackground: %v\nCreatedAt: %s\nStatus: %v\nStartTime: %s\nEndTime: %v\nElapsedTime: %v\nIsCompleted: %v\n\n",
			a.Id(), a.Transaction().Id(), a.Type(), a.Category(), a.Group(), a.ParentActionId(), a.IsBackground(), a.CreatedAt(), a.Status(), a.StartTime(), a.EndTime(), a.ElapsedTime(), a.IsCompleted())
	}

	fmt.Printf("ActionManager:\nCounter: %d\nNumCreated: %d\nNumInProgress: %d\n", actionManager.Counter(), actionManager.NumCreated(), actionManager.NumInProgress())
}

func testOperationManager(tranManager *actions.TransactionManager, actionManager *actions.ActionManager) {
	fmt.Println("***** testOperationManager *****")
	t, err := tranManager.CreateAndStart()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Transaction:\nId: %s\nCreatedAt: %s\nStartTime: %s\n\n", t.Id(), t.CreatedAt(), t.StartTime())

	a, err := actionManager.CreateAndStart(t, actions.ActionTypeApplication_Start, actions.ActionCategoryCommon, actions.ActionGroupApplication, uuid.NullUUID{}, false)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Action:\nId: %s\nTranId: %s\nType: %v\nCategory: %v\nGroup: %v\nParentActionId: %v\nIsBackground: %v\nCreatedAt: %s\nStatus: %v\nStartTime: %s\nEndTime: %v\nElapsedTime: %v\nIsCompleted: %v\n\n",
		a.Id(), a.Transaction().Id(), a.Type(), a.Category(), a.Group(), a.ParentActionId(), a.IsBackground(), a.CreatedAt(), a.Status(), a.StartTime(), a.EndTime(), a.ElapsedTime(), a.IsCompleted())

	for i := 1; i <= 10; i++ {
		o, err := a.Operations.CreateAndStart(actions.OperationTypeApplication_Start, actions.OperationCategoryCommon, actions.OperationGroupApplication, uuid.NullUUID{},
			actions.NewOperationParam("key1", "value1"),
			actions.NewOperationParam("key2", 2),
			actions.NewOperationParam("key3", "value3"),
		)

		if err != nil {
			panic(err)
		}

		fmt.Printf("Operation:\nId: %s\nTranId: %s\nActionId: %s\nType: %v\nCategory: %v\nGroup: %v\nParentOperationId: %v\nCreatedAt: %s\nStatus: %v\nStartTime: %s\nEndTime: %v\nElapsedTime: %v\nIsCompleted: %v\n",
			o.Id(), o.Action().Transaction().Id(), o.Action().Id(), o.Type(), o.Category(), o.Group(), o.ParentOperationId(), o.CreatedAt(), o.Status(), o.StartTime(), o.EndTime(), o.ElapsedTime(), o.IsCompleted())

		fmt.Println("Params:")

		for i, p := range o.Params() {
			fmt.Printf("[%d] Name: %s; Value: %v\n", i, p.Name, p.Value)
		}

		fmt.Println()

		if err = a.Operations.Complete(o, i%2 == 1); err != nil {
			panic(err)
		}

		fmt.Printf("Operation:\nId: %s\nTranId: %s\nActionId: %s\nType: %v\nCategory: %v\nGroup: %v\nParentOperationId: %v\nCreatedAt: %s\nStatus: %v\nStartTime: %s\nEndTime: %v\nElapsedTime: %v\nIsCompleted: %v\n",
			o.Id(), o.Action().Transaction().Id(), o.Action().Id(), o.Type(), o.Category(), o.Group(), o.ParentOperationId(), o.CreatedAt(), o.Status(), o.StartTime(), o.EndTime(), o.ElapsedTime(), o.IsCompleted())

		fmt.Println("Params:")

		for i, p := range o.Params() {
			fmt.Printf("[%d] Name: %s; Value: %v\n", i, p.Name, p.Value)
		}

		fmt.Println()
	}

	if err = actionManager.Complete(a, true); err != nil {
		panic(err)
	}

	fmt.Printf("Action:\nId: %s\nTranId: %s\nType: %v\nCategory: %v\nGroup: %v\nParentActionId: %v\nIsBackground: %v\nCreatedAt: %s\nStatus: %v\nStartTime: %s\nEndTime: %v\nElapsedTime: %v\nIsCompleted: %v\n\n",
		a.Id(), a.Transaction().Id(), a.Type(), a.Category(), a.Group(), a.ParentActionId(), a.IsBackground(), a.CreatedAt(), a.Status(), a.StartTime(), a.EndTime(), a.ElapsedTime(), a.IsCompleted())

	fmt.Printf("OperationManager:\nCounter: %d\nNumCreated: %d\nNumInProgress: %d\n", a.Operations.Counter(), a.Operations.NumCreated(), a.Operations.NumInProgress())
}

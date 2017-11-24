// Copyright 2017 The etcd-operator Authors
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

package e2eslow

import (
	"os"
	"testing"

	"github.com/coreos/etcd-operator/test/e2e/framework"

	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	if err := framework.Setup(); err != nil {
		logrus.Errorf("fail to setup framework: %v", err)
		os.Exit(1)
	}

	code := m.Run()

	if err := framework.Teardown(); err != nil {
		logrus.Errorf("fail to teardown framework: %v", err)
		os.Exit(1)
	}
	os.Exit(code)
}

package dandelion_facade

import (
	"genesis/pkg/core/dandelion/mixture"
	"genesis/pkg/util"
)

func TempalteGenerate(dto TemplateGenerateDTO) error {
	var qd mixture.ManagerStartOp

	if err := util.Copy(&qd, &dto); err != nil {
		return err
	}

	if err := mixture.NewManager(&qd).Start(); err != nil {
		return err
	}

	return nil
}

// Code generated by pegomock. DO NOT EDIT.
package matchers

import (
	"reflect"
	"github.com/petergtz/pegomock"
	gits "github.com/jenkins-x/jx/pkg/gits"
)

func AnySliceOfPtrToGitsGitRepository() []*gits.GitRepository {
	pegomock.RegisterMatcher(pegomock.NewAnyMatcher(reflect.TypeOf((*([]*gits.GitRepository))(nil)).Elem()))
	var nullValue []*gits.GitRepository
	return nullValue
}

func EqSliceOfPtrToGitsGitRepository(value []*gits.GitRepository) []*gits.GitRepository {
	pegomock.RegisterMatcher(&pegomock.EqMatcher{Value: value})
	var nullValue []*gits.GitRepository
	return nullValue
}

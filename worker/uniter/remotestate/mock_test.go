// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package remotestate_test

import (
	"github.com/juju/juju/api/watcher"
	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/state/multiwatcher"
	"github.com/juju/juju/worker/uniter/remotestate"
	"github.com/juju/names"
	"gopkg.in/juju/charm.v5"
)

type mockNotifyWatcher struct {
	changes chan struct{}
	stopped bool
}

func (w *mockNotifyWatcher) Stop() error {
	w.stopped = true
	return nil
}

func (*mockNotifyWatcher) Err() error {
	return nil
}

func (w *mockNotifyWatcher) Changes() <-chan struct{} {
	return w.changes
}

type mockStringsWatcher struct {
	changes chan []string
	stopped bool
}

func (w *mockStringsWatcher) Stop() error {
	w.stopped = true
	return nil
}

func (*mockStringsWatcher) Err() error {
	return nil
}

func (w *mockStringsWatcher) Changes() <-chan []string {
	return w.changes
}

type mockRelationUnitsWatcher struct {
	changes chan multiwatcher.RelationUnitsChange
	stopped bool
}

func (w *mockRelationUnitsWatcher) Stop() error {
	w.stopped = true
	return nil
}

func (*mockRelationUnitsWatcher) Err() error {
	return nil
}

func (w *mockRelationUnitsWatcher) Changes() <-chan multiwatcher.RelationUnitsChange {
	return w.changes
}

type mockState struct {
	unit                  mockUnit
	relations             map[names.RelationTag]*mockRelation
	storageAttachmentLife map[params.StorageAttachmentId]params.Life
	relationUnitsWatchers map[names.RelationTag]*mockRelationUnitsWatcher
}

func (st *mockState) Relation(tag names.RelationTag) (remotestate.Relation, error) {
	r, ok := st.relations[tag]
	if !ok {
		return nil, &params.Error{Code: params.CodeNotFound}
	}
	return r, nil
}

func (st *mockState) StorageAttachmentLife(
	ids []params.StorageAttachmentId,
) ([]params.LifeResult, error) {
	results := make([]params.LifeResult, len(ids))
	for i, id := range ids {
		life, ok := st.storageAttachmentLife[id]
		if !ok {
			results[i] = params.LifeResult{
				Error: &params.Error{Code: params.CodeNotFound},
			}
			continue
		}
		results[i] = params.LifeResult{Life: life}
	}
	return results, nil
}

func (st *mockState) Unit(tag names.UnitTag) (remotestate.Unit, error) {
	if tag != st.unit.tag {
		return nil, &params.Error{Code: params.CodeNotFound}
	}
	return &st.unit, nil
}

func (st *mockState) WatchRelationUnits(
	relationTag names.RelationTag, unitTag names.UnitTag,
) (watcher.RelationUnitsWatcher, error) {
	if unitTag != st.unit.tag {
		return nil, &params.Error{Code: params.CodeNotFound}
	}
	watcher, ok := st.relationUnitsWatchers[relationTag]
	if !ok {
		return nil, &params.Error{Code: params.CodeNotFound}
	}
	return watcher, nil
}

type mockUnit struct {
	tag                   names.UnitTag
	life                  params.Life
	resolved              params.ResolvedMode
	service               mockService
	unitWatcher           mockNotifyWatcher
	addressesWatcher      mockNotifyWatcher
	configSettingsWatcher mockNotifyWatcher
	storageWatcher        mockStringsWatcher
}

func (u *mockUnit) Life() params.Life {
	return u.life
}

func (u *mockUnit) Refresh() error {
	return nil
}

func (u *mockUnit) Resolved() (params.ResolvedMode, error) {
	return u.resolved, nil
}

func (u *mockUnit) Service() (remotestate.Service, error) {
	return &u.service, nil
}

func (u *mockUnit) Tag() names.UnitTag {
	return u.tag
}

func (u *mockUnit) Watch() (watcher.NotifyWatcher, error) {
	return &u.unitWatcher, nil
}

func (u *mockUnit) WatchAddresses() (watcher.NotifyWatcher, error) {
	return &u.addressesWatcher, nil
}

func (u *mockUnit) WatchConfigSettings() (watcher.NotifyWatcher, error) {
	return &u.configSettingsWatcher, nil
}

func (u *mockUnit) WatchStorage() (watcher.StringsWatcher, error) {
	return &u.storageWatcher, nil
}

type mockService struct {
	tag                   names.ServiceTag
	life                  params.Life
	curl                  *charm.URL
	forceUpgrade          bool
	serviceWatcher        mockNotifyWatcher
	leaderSettingsWatcher mockNotifyWatcher
	relationsWatcher      mockStringsWatcher
}

func (s *mockService) CharmURL() (*charm.URL, bool, error) {
	return s.curl, s.forceUpgrade, nil
}

func (s *mockService) Life() params.Life {
	return s.life
}

func (s *mockService) Refresh() error {
	return nil
}

func (s *mockService) Tag() names.ServiceTag {
	return s.tag
}

func (s *mockService) Watch() (watcher.NotifyWatcher, error) {
	return &s.serviceWatcher, nil
}

func (s *mockService) WatchLeadershipSettings() (watcher.NotifyWatcher, error) {
	return &s.leaderSettingsWatcher, nil
}

func (s *mockService) WatchRelations() (watcher.StringsWatcher, error) {
	return &s.relationsWatcher, nil
}

type mockRelation struct {
	id   int
	life params.Life
}

func (r *mockRelation) Id() int {
	return r.id
}

func (r *mockRelation) Life() params.Life {
	return r.life
}

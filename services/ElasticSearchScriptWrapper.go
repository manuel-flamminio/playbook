package services

import (
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"github.com/google/uuid"
)

type ScriptWrapper struct {
	script *types.Script
	source string
	params map[string]string
}

func NewScriptWrapper() *ScriptWrapper {
	source := `
		if (ctx._source.starredByUser == null) {
			ctx._source.starredByUser = []
		}
		if (ctx._source.upvotedByUser == null) {
			ctx._source.upvotedByUser = []
		}
		if (ctx._source.downvotedByUser == null) {
			ctx._source.downvotedByUser = []
		}
		if (ctx._source.numberOfFailures == null) {
			ctx._source.numberOfFailures = 0
		}
		if (ctx._source.numberOfSuccesses == null) {
			ctx._source.numberOfSuccesses = 0
		}
		if (ctx._source.numberOfTries == null) {
			ctx._source.numberOfTries = 0
		}
	`
	return &ScriptWrapper{script: &types.Script{}, source: source, params: make(map[string]string, 3)}
}

func (s *ScriptWrapper) AddStarredByUser(userId uuid.UUID) *ScriptWrapper {
	s.params["starredByUser"] = userId.String()
	s.source += "ctx._source.starredByUser.add(params.starredByUser);"
	return s
}

func (s *ScriptWrapper) RemoveStarredByUser(userId uuid.UUID) *ScriptWrapper {
	s.params["starredByUser"] = userId.String()
	s.source += `
		for (int i = ctx._source.starredByUser.length - 1; i >= 0; i--) {
			if (ctx._source.starredByUser[i] == params.starredByUser) {
				ctx._source.starredByUser.remove(i);
			}
		}
	`
	return s
}

func (s *ScriptWrapper) RemoveDownvotedByUser(userId uuid.UUID) *ScriptWrapper {
	s.params["downvotedByUser"] = userId.String()
	s.source += `
		for (int i = ctx._source.downvotedByUser.length - 1; i >= 0; i--) {
			if (ctx._source.downvotedByUser[i] == params.downvotedByUser) {
				ctx._source.downvotedByUser.remove(i);
				ctx._source.numberOfFailures -= 1;
				ctx._source.numberOfTries -= 1;
				if (ctx._source.numberOfTries != 0) {
					ctx._source.successPercentage = (double) Math.round((double)ctx._source.numberOfSuccesses/(double)ctx._source.numberOfTries*100)/100
				} else {
					ctx._source.successPercentage = 0
				}
			}
		}
	`
	return s
}

func (s *ScriptWrapper) RemoveUpvotedByUser(userId uuid.UUID) *ScriptWrapper {
	s.params["upvotedByUser"] = userId.String()
	s.source += `
		for (int i = ctx._source.upvotedByUser.length - 1; i >= 0; i--) {
			if (ctx._source.upvotedByUser[i] == params.upvotedByUser) {
				ctx._source.upvotedByUser.remove(i);
				ctx._source.numberOfSuccesses -= 1; 
				ctx._source.numberOfTries -= 1; 
				if (ctx._source.numberOfTries != 0) {
					ctx._source.successPercentage = (double) Math.round((double)ctx._source.numberOfSuccesses/(double)ctx._source.numberOfTries*100)/100
				} else {
					ctx._source.successPercentage = 0
				}
			}
		}
	`
	return s
}

func (s *ScriptWrapper) AddDownvotedByUser(userId uuid.UUID) *ScriptWrapper {
	s.params["downvotedByUser"] = userId.String()
	s.source += `
		ctx._source.downvotedByUser.add(params.downvotedByUser);
		ctx._source.numberOfFailures += 1; 
		ctx._source.numberOfTries += 1; 
		if (ctx._source.numberOfTries != 0) {
			ctx._source.successPercentage = (double) Math.round((double)ctx._source.numberOfSuccesses/(double)ctx._source.numberOfTries*100)/100
		} else {
			ctx._source.successPercentage = 0
		}
	`
	return s
}

func (s *ScriptWrapper) AddUpvotedByUser(userId uuid.UUID) *ScriptWrapper {
	s.params["upvotedByUser"] = userId.String()
	s.source += `
		ctx._source.upvotedByUser.add(params.upvotedByUser);
		ctx._source.numberOfSuccesses += 1; 
		ctx._source.numberOfTries += 1; 
		if (ctx._source.numberOfTries != 0) {
			ctx._source.successPercentage = (double) Math.round((double)ctx._source.numberOfSuccesses/(double)ctx._source.numberOfTries*100)/100
		} else {
			ctx._source.successPercentage = 0
		}
	`
	return s
}

func (s *ScriptWrapper) GetScript() (*types.Script, error) {
	scriptParams := make(map[string]json.RawMessage, 3)
	for key, param := range s.params {
		jsonParam, err := json.Marshal(param)
		if err != nil {
			return nil, err
		}
		scriptParams[key] = jsonParam
	}

	s.script.Lang = &scriptlanguage.Painless
	s.script.Params = scriptParams
	s.script.Source = &s.source
	return s.script, nil
}

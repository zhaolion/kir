package config

import "os"

func (suite *Suite) TestAddConfigPath() {
	cfg := New()
	count := len(cfg.configPaths)
	cfg.AddPath("/test/path")
	suite.Equal(1, len(cfg.configPaths)-count, "should add a path")
}

func (suite *Suite) TestJSON() {
	cfg := &stubConfig{}
	gotErr := Initialize(suite.JSONFile, cfg)
	suite.NoError(gotErr)
	got, ok := Get().(*stubConfig)
	suite.True(ok)
	suite.Equal("test-app", got.AppName)
	suite.Equal("test", got.AppEnv)
	suite.Equal(true, got.Debug)
	suite.Equal("root@tcp(localhost:3306)/test", got.Database.DSN)
}

func (suite *Suite) TestENV() {
	_ = os.Setenv("APP_ENV", "development")
	_ = os.Setenv("APP_DEBUG", "false")
	_ = os.Setenv("DB_DSN", "test@tcp(localhost:3306)")

	cfg := &stubConfig{}
	gotErr := Initialize(suite.EnvFile, cfg)
	suite.NoError(gotErr)
	got, ok := Get().(*stubConfig)
	suite.True(ok)
	suite.Equal("test-app", got.AppName)
	suite.Equal("development", got.AppEnv)
	suite.Equal(false, got.Debug)
	suite.Equal("test@tcp(localhost:3306)", got.Database.DSN)
}

func (suite *Suite) TestYAML() {
	cfg := &stubConfig{}
	gotErr := Initialize(suite.YAMLFile, cfg)
	suite.NoError(gotErr)
	got, ok := Get().(*stubConfig)
	suite.True(ok)
	suite.Equal("test-app", got.AppName)
	suite.Equal("test", got.AppEnv)
	suite.Equal(true, got.Debug)
	suite.Equal("root@tcp(localhost:3306)/test", got.Database.DSN)
}

func (suite *Suite) TestTOML() {
	cfg := &stubConfig{}
	gotErr := Initialize(suite.TOMLFile, cfg)
	suite.NoError(gotErr)
	got, ok := Get().(*stubConfig)
	suite.True(ok)
	suite.Equal("test-app", got.AppName)
	suite.Equal("test", got.AppEnv)
	suite.Equal(true, got.Debug)
	suite.Equal("root@tcp(localhost:3306)/test", got.Database.DSN)
}

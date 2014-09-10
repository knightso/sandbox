'use strict';
var util = require('util');
var path = require('path');
var yeoman = require('yeoman-generator');
var chalk = require('chalk');
var EventEmitter = require('events').EventEmitter;

var spreadsheet = require('./spreadsheet/spreadsheet.js');
var cellflow = require('./spreadsheet/cellflow.js');
var spreadsheetUtil = require('./spreadsheet/util.js');

var client = new spreadsheet();
var ssutil = new spreadsheetUtil();

var MockGenerator = yeoman.generators.Base.extend({
  init: function () {
    this.pkg = yeoman.file.readJSON(path.join(__dirname, '../package.json'));

    this.on('end', function () {
      if (!this.options['skip-install']) {
        this.npmInstall();
      }
    });
  },

  hello: function() {
    // have Yeoman greet the user
    console.log(this.yeoman);

    // replace it with a short and sweet description of your generator
    console.log(chalk.magenta('AngularJSのMockファイルを生成するジェネレータです。'));
    //console.log(chalk.magenta('Google Spreadsheetをjsonで出力します。'));
  },

  token: function() {
    var done = this.async();
    var self = this;
    ssutil.readToken(function(err, data) {
      console.log('token err=' + err);
      console.log('token data=' + data);
      if(err) {
        console.log('Visit the URL: ', client.generateAuthUrl());
        var prompts = [{
          type: 'input',
          name: 'requestCode',
          message: 'Enter the code here.'
        }];

        self.prompt(prompts, function(props) {
          client.getAccessToken(props.requestCode, function(tokens) {
            ssutil.saveToken(tokens, function(err) {
              if(err) {
                console.error(err);
              }
              done();
            });
          });
        });
      }else {
        client.oAuth2Client.setCredentials(JSON.parse(data.toString()));
        done();
      }
    });
  },

  askFor: function () {
    var done = this.async();
    var prompts = [{
      type: 'input',
      name: 'spreadsheetKey',
      message: 'Enter the spreadsheet key here.'
    }];

    this.prompt(prompts, function (props) {
      this.spreadsheetKey = props.spreadsheetKey;

      done();
    }.bind(this));
  },

  table: function() {
    var done = this.async();
    var ev = new EventEmitter;
    var self = this;

    var flow = new cellflow(client);
    ev.once('return_data', function(data) {
      self.write(data.title+'.json', JSON.stringify(data, function(key, value) {
        if (value instanceof Object === false || Object.getPrototypeOf(value) !== Object.prototype) {
          return value;
        }
        var keys = Object.keys(value);
        keys.sort();
        var newValue = {};
        keys.forEach(function(key) {
          newValue[key] = value[key];
        });
        return newValue;
      }, 2));
      done();
    });

    flow.main(this.spreadsheetKey, ev);
  },

  app: function () {
    // mock用のファイルをコピーする。

    this.copy('_package.json', 'package.json');
    this.copy('_bower.json', 'bower.json');
  },

  projectfiles: function () {
    this.copy('editorconfig', '.editorconfig');
    this.copy('jshintrc', '.jshintrc');
  },
});

module.exports = MockGenerator;

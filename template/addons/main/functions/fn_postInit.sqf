addMissionEventHandler ["ExtensionCallback", {
  params ["_extension", "_function", "_args"];

  if !(_extension isEqualTo "EXTENSION_NAME") exitWith {};

  _argsArr = parseSimpleArray _args;
  if (count _argsArr isEqualTo 0) exitWith {
    diag_log format["a3go: No arguments received from extension. %1", _args];
  };

  switch (_function) do {
    case "testAsync": {
      diag_log format["a3go: ""testAsync"" callback received from extension. %1", _argsArr];
    };
    case "test": {
      diag_log format["a3go: ""test"" callback received from extension. %1", _argsArr];
    };
    case "saveMyCall": {
      diag_log format["a3go: ""saveMyCall"" callback received from extension. %1", _argsArr];
    };
	default {
	  diag_log format["a3go: Unknown function %1", _function];
	};
  };
}];
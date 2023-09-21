class CfgPatches {
	class TestAddon {
		units[] = {};
		weapons[] = {};
		requiredVersion = 0.1;
		requiredAddons[] = {"A3_Data_F"};
	};
};

class CfgFunctions {
	class testa {
		class functions {
			file = "x\addons\testa\main\functions";
			class test {};
		};
	};
};
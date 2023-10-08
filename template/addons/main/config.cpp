class CfgPatches {
	class a3go_main {
		units[] = {};
		weapons[] = {};
		requiredVersion = 0.1;
		requiredAddons[] = {"A3_Data_F"};
	};
};

class CfgFunctions {
	class a3go {
		class functions {
			file = "x\addons\a3go\main\functions";
			class postInit {postInit = 1;};
			class testSync {};
			class testAsync {};
			class testSaveCaller {};
			class hashToJson {};
		};
	};
};
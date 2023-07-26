package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	RESOURCE_PROVIDERS_STR = "--RESOURCEPROVIDERS--"
	CLUSTER_PROVIDERS_STR  = "--CLUSTERPROVIDERS--"
	TEAM_STR               = "--TEAMS--"
	ADDONS_STR             = "--ADDONS--"

	start = "Add"
	end   = "Request"
)

var (
	objectFileStr = flag.String("obj", "obj.txt", "path to object text file")
	addonFileStr  = flag.String("addon", "addons.proto", "path to addon proto file")
	clPFileStr    = flag.String("clp", "cluster_providers.proto", "path to cluster provider proto file")
	rPFileStr     = flag.String("rp", "resource_providers.proto", "path to resource provider proto file")
	teamFileStr   = flag.String("team", "teams.proto", "path to team proto file")
	rpcFileStr    = flag.String("rpc", "cluster.proto", "path to proto file that defines rpc service")
)

func main() {
	flag.Parse()
	file, err := os.Open(*objectFileStr)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	clusterProviders := make([]string, 0, 5)
	teams := make([]string, 0, 5)
	resourceProviders := make([]string, 0, 5)
	addons := make([]string, 0, 10)
	// arr := ""
	for scanner.Scan() {
		line := scanner.Text()
		//if line == ADDONS_STR || line == CLUSTER_PROVIDERS_STR || line == RESOURCE_PROVIDERS_STR || line == TEAM_STR {
		//	arr = line
		//	continue
		//}

		// if arr == ADDONS_STR {
		// 	addons = append(addons, line)
		// } else if arr == CLUSTER_PROVIDERS_STR {
		// 	clusterProviders = append(clusterProviders, line)
		// } else if arr == RESOURCE_PROVIDERS_STR {
		// 	resourceProviders = append(resourceProviders, line)
		// } else if arr == TEAM_STR {
		// 	teams = append(teams, line)
		// }
		if strings.Contains(line, "AddOn") {
			addons = append(addons, line)
		} else if strings.Contains(line, "ClusterProvider") {
			clusterProviders = append(clusterProviders, line)
		} else if strings.Contains(line, "Provider") {
			resourceProviders = append(resourceProviders, line)
		} else if strings.Contains(line, "Team") {
			teams = append(teams, line)
		}

	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	file.Close()
	addonF, err := os.OpenFile(*addonFileStr, os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	clusterProvidersF, err := os.OpenFile(*clPFileStr, os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	resourceProvidersF, err := os.OpenFile(*rPFileStr, os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	teamF, err := os.OpenFile(*teamFileStr, os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	checkAndWrite(addonF, addons)
	checkAndWrite(resourceProvidersF, resourceProviders)
	checkAndWrite(clusterProvidersF, clusterProviders)
	checkAndWrite(teamF, teams)

	rpcsRead, err := os.OpenFile(*rpcFileStr, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	rpcsWrite, err := os.OpenFile(*rpcFileStr+".temp", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	writeRPCs(rpcsRead, rpcsWrite, append(addons, append(clusterProviders, append(resourceProviders, teams...)...)...))
}

func checkAndWrite(fileToWrite *os.File, objects []string) error {
	fileB, err := ioutil.ReadAll(fileToWrite)
	if err != nil {
		return err
	}
	fileStr := string(fileB)
	for _, obj := range objects {
		if !strings.Contains(fileStr, " "+obj) {
			err := writeObject(*fileToWrite, obj)
			if err != nil {
				return err
			}
		}
		if !strings.Contains(fileStr, start+obj+end) {
			err := writeAddRequest(*fileToWrite, obj)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func writeAddRequest(file os.File, object string) error {
	_, err := file.WriteString("\nmessage " + start + object + end + " {\n")
	if err != nil {
		return err
	}
	_, err = file.WriteString("\tstring cluster_name = 1;\n")
	if err != nil {
		return err
	}
	_, err = file.WriteString("\t" + object + " " + strings.ToLower(object) + " = 2;\n")
	if err != nil {
		return err
	}
	_, err = file.WriteString("}\n")
	if err != nil {
		return err
	}
	return nil
}
func writeRPCs(fileToRead *os.File, fileToWrite *os.File, objects []string) (string, error) {
	defer fileToRead.Close()
	defer fileToWrite.Close()
	fileB, err := ioutil.ReadAll(fileToWrite)
	if err != nil {
		return "", err
	}
	fileStr := string(fileB)
	objsToAdd := make([]string, 0, 10)
	for _, obj := range objects {
		if !strings.Contains(fileStr, "rpc "+start+obj+" ") {
			objsToAdd = append(objsToAdd, obj)
		}
	}
	scanner := bufio.NewScanner(fileToRead)
	checkNext := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "rpc") {
			checkNext = true
		}
		if checkNext && strings.Contains(line, "}") {
			for _, obj := range objsToAdd {
				_, err := fileToWrite.WriteString("\n\trpc " + start + obj + " (" + start + obj + end + ") returns (APIResponse);\n")
				if err != nil {
					return "", err
				}
			}
			checkNext = false

		}
		_, err := fileToWrite.WriteString(line + "\n")
		if err != nil {
			return "", err
		}
	}

	return "", nil
}

func writeObject(file os.File, object string) error {
	_, err := file.WriteString("\nmessage " + object + " {\n")
	if err != nil {
		return err
	}
	_, err = file.WriteString("\tstring temp = 1;\n")
	if err != nil {
		return err
	}
	_, err = file.WriteString("}\n")
	if err != nil {
		return err
	}
	return nil
}

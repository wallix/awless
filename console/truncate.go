/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package console

const truncateChars = "..."

func truncateLeft(str string, maxSize int) string {
	ltc := len(truncateChars)
	lstr := len(str)
	if maxSize < ltc {
		return str[lstr-maxSize : lstr]
	}
	if lstr > maxSize {
		return truncateChars + str[lstr-maxSize+ltc:lstr]
	}
	return str
}

func truncateRight(str string, maxSize int) string {
	ltc := len(truncateChars)
	lstr := len(str)
	if maxSize < ltc {
		return str[:maxSize]
	}
	if lstr > maxSize {
		return str[:maxSize-ltc] + truncateChars
	}
	return str
}

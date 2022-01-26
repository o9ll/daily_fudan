/**
 * @Author Oliver
 * @Date 1/26/22
 **/

package util

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

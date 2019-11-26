package hflag

import (
	"fmt"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHFlag(t *testing.T) {
	Convey("test flag", t, func() {
		flagSet := NewFlagSet("test flag")
		i := flagSet.Int("i", 10, "int flag")
		//i := new(int)
		//flagSet.IntVar(i, "i", 10, "int flag")
		f := flagSet.Float64("f", 11.11, "float flag")
		s := flagSet.String("s", "hello world", "string flag")
		d := flagSet.Duration("d", time.Duration(30)*time.Second, "string flag")
		b := flagSet.Bool("b", false, "bool flag")
		vi := flagSet.IntSlice("vi", []int{1, 2, 3, 4, 5}, "int slice flag")

		Convey("check default Value", func() {
			So(*i, ShouldEqual, 10)
			So(*f, ShouldAlmostEqual, 11.11)
			So(*s, ShouldEqual, "hello world")
			So(*d, ShouldEqual, time.Duration(30)*time.Second)
			So(*b, ShouldBeFalse)
			So(*vi, ShouldResemble, []int{1, 2, 3, 4, 5})
		})

		Convey("parse case success", func() {
			err := flagSet.Parse(strings.Split("-b -i 100 -f=12.12 --s golang --d=20s -vi 6,7,8,9", " "))
			So(err, ShouldBeNil)
			So(*i, ShouldEqual, 100)
			So(*f, ShouldAlmostEqual, 12.12)
			So(*s, ShouldEqual, "golang")
			So(*d, ShouldEqual, time.Duration(20)*time.Second)
			So(*b, ShouldBeTrue)
			So(*vi, ShouldResemble, []int{6, 7, 8, 9})
		})

		Convey("parse case unexpected Value", func() {
			// -i 后面期望之后一个参数，100 会被当做 -i 的参数，101 会被当成位置参数，继续解析
			err := flagSet.Parse(strings.Split("-b -i 100 101 -f 12.12 -s golang -d 20s", " "))
			So(err, ShouldBeNil)
			So(*b, ShouldBeTrue)
			So(*i, ShouldEqual, 100)
			So(*f, ShouldEqual, 12.12)
			So(*s, ShouldEqual, "golang")
			So(*d, ShouldEqual, 20*time.Second)
			So(flagSet.NFlag(), ShouldEqual, 6)
			So(flagSet.NArg(), ShouldEqual, 1)
			So(flagSet.Args(), ShouldResemble, []string{
				"101",
			})
		})

		Convey("set case", func() {
			err := flagSet.Set("i", "120")
			So(err, ShouldBeNil)
			So(*i, ShouldEqual, 120)
		})

		Convey("visit case", func() {
			// 遍历所有设置过得选项
			flagSet.Visit(func(f *Flag) {
				fmt.Println(f.Name)
			})

			// 遍历所有选项
			flagSet.VisitAll(func(f *Flag) {
				fmt.Println(f.Name)
			})
		})

		Convey("lookup case", func() {
			f := flagSet.Lookup("i")
			So(f.Name, ShouldEqual, "i")
			So(f.DefValue, ShouldEqual, "10")
			So(f.Usage, ShouldEqual, "int flag")
			So(f.Value.String(), ShouldEqual, "10")
		})

		Convey("print defaults", func() {
			flagSet.PrintDefaults()
		})
	})
}

func TestHFlagParse(t *testing.T) {
	Convey("test case1", t, func() {
		flagSet := NewFlagSet("test flag")
		So(flagSet.AddFlag("int-option", "i", "usage", "int", true, "0"), ShouldBeNil)
		So(flagSet.AddFlag("str-option", "s", "usage", "string", true, ""), ShouldBeNil)
		So(flagSet.AddFlag("key", "k", "usage", "float", true, ""), ShouldBeNil)
		So(flagSet.AddFlag("all", "a", "usage", "bool", true, ""), ShouldBeNil)
		So(flagSet.AddFlag("user", "u", "usage", "bool", true, ""), ShouldBeNil)
		So(flagSet.AddFlag("password", "p", "usage", "string", false, "654321"), ShouldBeNil)
		So(flagSet.AddPosFlag("pos1", "usage", "string", ""), ShouldBeNil)
		So(flagSet.AddPosFlag("pos2", "usage", "string", ""), ShouldBeNil)
		err := flagSet.Parse([]string{
			"pos1",
			"--int-option=123",
			"--str-option", "hello world",
			"-k", "3.14",
			"-au",
			"-p123456",
			"pos2",
		})
		So(err, ShouldBeNil)
		flagSet.Usage()
	})

	Convey("test case2", t, func() {
		flagSet := NewFlagSet("test flag")
		version := flagSet.Bool("v", false, "print current version")
		configfile := flagSet.String("c", "configs/monitor.json", "config file path")
		So(flagSet.Parse(strings.Split("--v", " ")), ShouldBeNil)
		So(*version, ShouldBeTrue)
		So(*configfile, ShouldEqual, "configs/monitor.json")
	})
}

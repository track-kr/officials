package main

import (
	"context"
	"os"

	"github.com/jszwec/csvutil"
	datagokr "github.com/track-kr/go-datagokr"
	"github.com/track-kr/go-datagokr/nationalassemblyinfoservice"
)

type AssemblyMember struct {
	nationalassemblyinfoservice.GetMemberCurrStateListResponseBodyItemsItem
	PolyNm string `csv:"polyNm"`
}

func main() {
	var items []nationalassemblyinfoservice.GetMemberCurrStateListResponseBodyItemsItem
	ctx := context.Background()
	for i := 1; i < 2; i++ {
		out, err := datagokr.DefaultClient.NationalAssemblyInfoService.GetMemberCurrStateList(ctx, &nationalassemblyinfoservice.GetMemberCurrStateListRequest{
			PageNo:    i,
			NumOfRows: 300,
		})
		if err != nil {
			panic(err)
		}
		items = append(items, out.Body.Items.Item...)
	}

	var assemblyMembers []AssemblyMember
	for i := 0; i < len(items); i++ {
		assemblyMembers = append(assemblyMembers, AssemblyMember{GetMemberCurrStateListResponseBodyItemsItem: items[i]})
	}

	for i := 0; i < len(items); i++ {
		infoout, err := datagokr.DefaultClient.NationalAssemblyInfoService.GetMemberDetailInfoList(ctx, &nationalassemblyinfoservice.GetMemberDetailInfoListRequest{
			DeptCd: items[i].DeptCd,
			EmpNm:  items[i].Num,
		})
		if err != nil {
			panic(err)
		}
		assemblyMembers[i].PolyNm = infoout.Body.Item.PolyNm
	}

	f, err := os.Create("../assembly.csv")
	if err != nil {
		panic(err)
	}

	b, err := csvutil.Marshal(assemblyMembers)
	if err != nil {
		panic(err)
	}
	f.Write(b)
}

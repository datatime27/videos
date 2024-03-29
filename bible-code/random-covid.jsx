
function addColorText(layer, index, color) {
    var animator = layer.property('Text').property('Animators').addProperty('ADBE Text Animator');
    var colorProp = animator.property('Properties').addProperty('ADBE Text Fill Color');
    colorProp.setValue(color);
    var selector = animator.property('Selectors').addProperty('ADBE Text Selector');
    var advanced = selector.property('Advanced');
    advanced.property('Based On').setValue(2); // Characters
    advanced.property('Units').setValue(2); // Index
    var start = selector.property('ADBE Text Index Start').setValue(index)
    var end = selector.property('ADBE Text Index End').setValue(index+1)
}


var font_size = 30;
var tracking = 700;

var positionX = 960;
var positionY = 540;

var width = 1800;
var height = 900;

var letterWidth = 22.2
var letterHeight = 36.0
var wrapLength = 22

var originX = positionX - width/2
var originY = positionY - height/2

var text = "MOWNEGEARECSWYHIHOUAEOIODHSSIGNRTODESIGSTSEPRBSLTNEPUICOANPNEITOTRNIAVTTTPTGLEBNLRESACTSEILEFOSILLLICTSIHIOAEEPKSGTNCTCRAANOTYENOTNOEUSNRAUAETROFCOTOYSLOALOUANTEEONCGESTNRSCSHEIRPAEOARLTRUUREPRMIKHAODGREAUSCODEOTPTELOENSORENDSLORANSUAHHTNUOEMSMRYAHLOOSNTSTADAUSEPYROIRUESRPUNDIEOUTGERENANATPAIIERENTETAWRACRHONLANOSNCTOGIRRUSATORBRTPCTCHEOEPVCISFEIILLIARMERMILTASEIPTPIPAALKLHDRLTAYAEACRAAUONNAERLSMAEELSCDEEAAUSHDMLWOVEEFDDUSAONMOAARARRAUSIAAEOOHOGUEEORRVDFEUTGBCELCMAORYMSRFOIHAEIADCROVIWEIONMOERSAEATMEARNDORSLIISAANTPRTSFSTAECOEDTNVCROEPNNAILUOYSNCOREUBLSCPVEOUDUACPNIEIEENIRORRRTSAOPSNRAEDUNASHEYOINPEFBVEARTEDUSCUIDLUCRNNNEPPISLANLAEESIIESCEIHEEGIAHSMIDLOWAAUCPDMIPNTGODPETANSESIMNBGNOPEHKATDTYTBNESFSAMNLRNIENOSUHADCEAONTMIPAIAOTESNEAREONAAPLAEKSIEYESSOSRHOLTRNYROCECANTTEUTGROITBFMENRIEIORNPHUSWANOTLEETSITTRLIEEICLEETRECGESCTCDFLHEIDURANNSITNBANLALOKTLSRBHCYAMRRIIOETHPTASEABBILTINRLHAENIIOAUOIRGTASVBNNRRIEPOHHIAYAOERTAEDDPNPESIESDBRWYNTSLTSBIMNOLRNOYEAESSOFILIUAEOHIEILRBLLCEEIDUEESTSFNWEEPNSCENEIHAHEACDCFLEEMVSIRDMSEADRPAELDNSNSNLESBSLVOAKAOHIAPHTEERIEMTEWAEECUPEMHRLOCNINTHAHSIETINRTSWESOADROOASOBECSIAANLRIYOOMOYTEESETIAONSAEDLURSRRRINREOREIDEINMOOWEAHNISRSAEEEELYINRREGLYDOKIGENYYDETMTNHASWAEIIAREFIANESREIMATTNIOGRAECPLRHTOOLAARLMMLHLIIDEOAKNDCMANAPEMRBRNAONURREALAAFMHHICICPITLAUEOIALTTKNIHEIUSTOBDGBALRGTURNWATHRENGHEHNRRTGPMENERCISOMHBRCAFAONROETETSTOPAIEEEEOESENEAARRRNTAHRRCRLADAAANUTINOCAOYTNEDTIPYESAMRTFICDMREALEPESUESILRSOAOHNUARGHLTNNNIRALGHRIAOOFOKNNTLRSERIFULLORPSAAETMIOCECPLRAANSIBLMLUUAAIIRNCELHAENNTIABNHHORTKRRRTOISCITOEAOOUPININLEHLPINCWIAOAEDCONELTEEVELCNSECDNIMFAPCEOTFNUPARTMHLURTUCUDIULRNSYOOEUDTOYEHIDEEOIAGIAUCDGSFEIINNAERDRRGCELRCOPTIBLHIOSERWIEERIAHTIIUATORNAUNMEORTIPCORIHITCOFSHEESTSIOBLICRWBEABITEPPUROMLMCHCEMSOIHSFNERIORICGWGHPRGARELHEATPNEITINUCETKASETTSIOEENKUEPIOLPGNGUTLAGYOCTIMRLEEMICBEKDTARAAGCMTARNLNNDEOAAEODRMBCTCTRHEOTLHCBHIIOHEEIRNCAPIEKCNISTSAAOAEDNOGIHITGETAMLSAPLROITADSARPMGLOCPTKRRRRETOBNAEEASYIRTUSAACLROTTEHUDRRSIFCUHEACTREHRAUERRUSETRSOENSTNEUCRIAREOAATTLESELTCRAWTITUHARNAPLITANGEBULRLBSLPBIOTSDLEPISEVNAEOONGEARTUTMIETSPTACROROATDCCONUVOUMOLOUTSSPROSCODPINPUSSIAEARIRISNAHNECNWDRMUIRRAGNCLYYMEEOEFFUVNEPAIMOBLIMSIGTEHRRTKECREEPDANUETMEMRNMPETNMIOAOOHITEAWERELHREROGUEYWESEIIMOUTODMOHKGOELABPNNIDLCWTFIDBUICEUOCCLATVGLLAIGNEEHUCIDCIYHRRRNUOAEWLDCHAAARATAEYGTCFNUROMEHEBUCMCCSINNCRENAUIFGBDNADSNLALOWISDSRENRAMMHEASUSRYINMTUVPIREUGIIDEAIIAMAMTLTLIALUWHIECACFLRORMSTAOHFDOBNTOHAOTDWTPGOUREAETIASYCIKECLAOUSEYENABCOICAELMEFORORAGAECKNEGTTEMRRPENEEEOEASABDREAYDLEEIPHEETRARCOIMATMEHIOWNERAATRLFESLORNTMNOUIAARHNEMHSEHAONNIOPNETNSEOOIICMLFLIONRNITTRGTIITNHRTCGWNCOOEETOITRAEPNBEUPECSMRSNLASEISURLEIGIENENATWDSNEENEAIRUSRCSSCCPOALSHDRNUIOCOSAEGLDHLGIIOHAAVONAIFIBRRRSRAAHUKCHNOSALADSHLTECNNPIRYNREFKAABHNEGHAGAIEPNMIENENNAUOPTCGCRMOOSNRNDHFOGDSPOCEAPSCULNGLTNPDKERORMORSNFVNSUENAIMIRDUGSNASTBWOIITRYEOICPNINMILODLOWOETUEASEHNRNRERUIATAAYYMOAKOIACRDISENARDOLMYRPWNRTUETIDMEEUNSILCBLORRNEIILHSMECIRAEMUTGERMSORERDONSHNNGNNSATTREARDESBRAESAKHIFNNCOAKRANARMITOOEEIISTGWBAOYNOCMTSGKTSRIUUASNNAORRRTARETYHEERUTMPERSHRDPEAFSEAEHIOCNCIPTENFASLRLANUIHYANMRCOEWEENHDANSTDAANLIIOOTANHNTDISKOUDIPEDMSSSTCRUTDIONRDKILNSLTCSETDNOFOTPSENGESRBOIENRKILRACNGLUEONMRRIEAGEUIWOLAHPELAERTCIOEPARTAYNBENSIITPRUSIHCOACISRCNEYTLHSROEOCBVLTNCELTFLBNOANPLPEAITODSAOAUCOELPTSECMORISEOTIORMAEWOSEMNTIPELEPEATKGLEDUWEAANFAEETRTESEEOHHIAMASBWICEIIHNKICFNUHUBHETGEBESSDNESIUOIATIUEHEPGEOFECELVAIEWETEMEEOMRSLAUERSRORSGSTANKAEYEICESOESRBTABHOBACRASTIOAETETBVBBUYODILOETVCDESGNEUAPIARHRWSFLANLLPECOUEPBNOIOOSETAITHSDRENNIERDEHTUTLOYABRRDLEONUIIHAUREMCSUATETSSHFNLEDYNIHCUVETRBCVUAOFIAMSNNOSATGWCFIIIFATGDNBNGULFERLAHASIGSBDENDAAPSNNMUIDEEERDERTTIRGYOOULORNACUTTNPKIEALPSFANTFEULNESGHICNLCAENTLTPHBDRAAIBIGTERMMPYEEBSEAAENSEEIARTELNTAAPAROAGONEINSETPNCMALREAOIUELSLBAEDAAENOTEAPENHLNSGHONREODOGTSTLUEOAGNYLEIUMOUIHOOCOOGOFINIDYTPEYHLSUNEERWKFITPPIGRAVTSSTOUANLETVDBIRTBPTACNLNUPTCPIMNLUYEETAHEILDEIUTMNWAMCPRFOAHWEIEOLSNAII";

app.beginUndoGroup('bible-code');
var comp = app.project.activeItem;

label = comp.layers.addBoxText([width,height], text);
label.threeDLayer = true

var textProp = label.property("Source Text");
label.position.setValue([positionX, positionY]);
var textDocument = textProp.value;
textDocument.fontSize = font_size;
textDocument.fillColor = [0, 0, 0];
textDocument.tracking = tracking;
textProp.setValue(textDocument);

// "WUHAN": 532246-533174 skip:232 
addColorText(label, 2923, [0.9, 0.25760633, 0.0])  //W
addColorText(label, 3155, [0.9, 0.25760633, 0.0])  //U
addColorText(label, 3387, [0.9, 0.25760633, 0.0])  //H
addColorText(label, 3619, [0.9, 0.25760633, 0.0])  //A
addColorText(label, 3851, [0.9, 0.25760633, 0.0])  //N

// "SICK": 531436-532486 skip:350 
addColorText(label, 2113, [0.9, 0.0, 0.20318687])  //S
addColorText(label, 2463, [0.9, 0.0, 0.20318687])  //I
addColorText(label, 2813, [0.9, 0.0, 0.20318687])  //C
addColorText(label, 3163, [0.9, 0.0, 0.20318687])  //K


// "COVID": 531303-531679 skip:94 
addColorText(label, 1980, [0.0, 0.9, 0.9])  //C
addColorText(label, 2074, [0.0, 0.9, 0.9])  //O
addColorText(label, 2168, [0.0, 0.9, 0.9])  //V
addColorText(label, 2262, [0.0, 0.9, 0.9])  //I
addColorText(label, 2356, [0.0, 0.9, 0.9])  //D

/*
// "SURIV": 531121-532597 skip:369 
addColorText(label, 1798, [0.0, 0.9, 0])  //S
addColorText(label, 2167, [0.0, 0.9, 0])  //U
addColorText(label, 2536, [0.0, 0.9, 0])  //R
addColorText(label, 2905, [0.0, 0.9, 0])  //I
addColorText(label, 3274, [0.0, 0.9, 0])  //V

// "KCIS": 530017-530431 skip:138 
addColorText(label, 694, [0.0, 0.9, 0.09696952])  //K
addColorText(label, 832, [0.0, 0.9, 0.09696952])  //C
addColorText(label, 970, [0.0, 0.9, 0.09696952])  //I
addColorText(label, 1108, [0.0, 0.9, 0.09696952])  //S

// "SICK": 530270-532361 skip:697 
addColorText(label, 947, [0.0, 0.9, 0.6645601])  //S
addColorText(label, 1644, [0.0, 0.9, 0.6645601])  //I
addColorText(label, 2341, [0.0, 0.9, 0.6645601])  //C
addColorText(label, 3038, [0.0, 0.9, 0.6645601])  //K

// "SICK": 530762-531512 skip:250 
addColorText(label, 1439, [0.0, 0.9, 0.31805816])  //S
addColorText(label, 1689, [0.0, 0.9, 0.31805816])  //I
addColorText(label, 1939, [0.0, 0.9, 0.31805816])  //C
addColorText(label, 2189, [0.0, 0.9, 0.31805816])  //K

// "KCIS": 531119-532004 skip:295 
addColorText(label, 1796, [0.0, 0.9, 0.29310185])  //K
addColorText(label, 2091, [0.0, 0.9, 0.29310185])  //C
addColorText(label, 2386, [0.0, 0.9, 0.29310185])  //I
addColorText(label, 2681, [0.0, 0.9, 0.29310185])  //S



// "CHINA": 529423-531639 skip:554 
addColorText(label, 100, [0.0, 0.9, 0.752573])  //C
addColorText(label, 654, [0.0, 0.9, 0.752573])  //H
addColorText(label, 1208, [0.0, 0.9, 0.752573])  //I
addColorText(label, 1762, [0.0, 0.9, 0.752573])  //N
addColorText(label, 2316, [0.0, 0.9, 0.752573])  //A

// "SICK": 530029-531787 skip:586 
addColorText(label, 706, [0.0, 0.9, 0.69671917])  //S
addColorText(label, 1292, [0.0, 0.9, 0.69671917])  //I
addColorText(label, 1878, [0.0, 0.9, 0.69671917])  //C
addColorText(label, 2464, [0.0, 0.9, 0.69671917])  //K

// "SICK": 530363-532322 skip:653 
addColorText(label, 1040, [0.0, 0.9, 0.15832828])  //S
addColorText(label, 1693, [0.0, 0.9, 0.15832828])  //I
addColorText(label, 2346, [0.0, 0.9, 0.15832828])  //C
addColorText(label, 2999, [0.0, 0.9, 0.15832828])  //K
*/
app.endUndoGroup();

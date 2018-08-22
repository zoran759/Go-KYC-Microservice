package common

type KYCResult int

const Error KYCResult = -1
const Approved KYCResult = 1
const Denied KYCResult = 2
const Unclear KYCResult = 3

type KYCFinality int

const Final KYCFinality = 1
const NonFinal KYCFinality = 2
const Unknown KYCFinality = 3

type Gender int

const Male Gender = 1
const Female Gender = 2
